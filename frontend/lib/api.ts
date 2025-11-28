/**
 * API Client for The Asheville Setlist
 *
 * Base fetch wrapper with error handling and type-safe response parsing
 */

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * Standard API response envelope
 */
export interface APIResponse<T> {
  data: T;
  meta?: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

/**
 * Standard API error response
 */
export interface APIError {
  error: {
    code: string;
    message: string;
    details?: Record<string, string[]>;
  };
}

/**
 * Custom error class for API errors
 */
export class APIRequestError extends Error {
  constructor(
    message: string,
    public code: string,
    public details?: Record<string, string[]>
  ) {
    super(message);
    this.name = 'APIRequestError';
  }
}

/**
 * Base fetch wrapper with error handling
 */
export async function apiRequest<T>(
  endpoint: string,
  options?: RequestInit
): Promise<APIResponse<T>> {
  const url = `${API_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      const error = data as APIError;
      throw new APIRequestError(
        error.error.message,
        error.error.code,
        error.error.details
      );
    }

    return data as APIResponse<T>;
  } catch (error) {
    if (error instanceof APIRequestError) {
      throw error;
    }

    // Network or parsing errors
    throw new APIRequestError(
      error instanceof Error ? error.message : 'An unexpected error occurred',
      'NETWORK_ERROR'
    );
  }
}

/**
 * GET request helper
 */
export async function apiGet<T>(
  endpoint: string,
  params?: URLSearchParams
): Promise<APIResponse<T>> {
  const url = params ? `${endpoint}?${params.toString()}` : endpoint;
  return apiRequest<T>(url);
}

/**
 * POST request helper
 */
export async function apiPost<T, D = unknown>(
  endpoint: string,
  data: D
): Promise<APIResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

/**
 * PUT request helper
 */
export async function apiPut<T, D = unknown>(
  endpoint: string,
  data: D
): Promise<APIResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

/**
 * DELETE request helper
 */
export async function apiDelete<T>(endpoint: string): Promise<APIResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: 'DELETE',
  });
}
