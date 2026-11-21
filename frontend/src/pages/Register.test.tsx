import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import Register from './Register'

function mockJson(body: unknown, status: number) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function fillAndSubmit(employeeId: string) {
  render(<Register />)
  fireEvent.change(screen.getByLabelText('Employee ID'), {
    target: { value: employeeId },
  })
  fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Amy' } })
  fireEvent.change(screen.getByLabelText('Password'), {
    target: { value: 'secret' },
  })
  fireEvent.click(screen.getByRole('button', { name: 'Register' }))
}

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('Register', () => {
  it('blocks a bad employee id without calling the API', async () => {
    const fetchMock = mockJson({}, 201)

    fillAndSubmit('123')

    expect(await screen.findByText('Employee ID must be 6 digits')).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('shows a 409 as an employee id field error', async () => {
    mockJson(
      {
        code: 'EMPLOYEE_ID_TAKEN',
        message: 'Already registered',
        requestId: 'r1',
      },
      409,
    )

    fillAndSubmit('123456')

    expect(await screen.findByText('Already registered')).toBeTruthy()
  })

  it('maps 422 fields onto their inputs', async () => {
    mockJson(
      {
        code: 'VALIDATION_FAILED',
        message: 'Invalid',
        requestId: 'r2',
        fields: [{ field: 'name', message: 'Name is too long' }],
      },
      422,
    )

    fillAndSubmit('123456')

    expect(await screen.findByText('Name is too long')).toBeTruthy()
  })

  it('shows a form error when fetch throws', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('offline')))

    fillAndSubmit('123456')

    expect(await screen.findByText(/Network error/)).toBeTruthy()
  })
})
