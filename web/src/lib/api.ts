const apiUrl = import.meta.env.VITE_PUBLIC_API_URL

// Generic API request handler
async function apiRequest(endpoint: string, options: RequestInit = {}) {
  const response = await fetch(`${apiUrl}${endpoint}`, {
    ...options,
    credentials: 'include',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
  })

  const data = await response.json()

  if (!response.ok) {
    return { error: data.error || 'Something went wrong' }
  }

  return { data }
}

export async function getLoggedInUser() {
  return await apiRequest('/admin/me', {
    method: 'GET',
  })
}

export async function fetchProducts() {
  console.log('Fetching products')
  return await apiRequest('/admin/products', {
    method: 'GET',
  })
}

export async function signInWithPassword({
  email,
  password,
}: {
  email: string
  password: string
}) {
  return await apiRequest('/admin/sign-in', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function listCustomers() {
  return await apiRequest('/admin/customers', {
    method: 'GET',
  })
}

export async function listOrders() {
  return await apiRequest('/admin/orders', {
    method: 'GET',
  })
}
