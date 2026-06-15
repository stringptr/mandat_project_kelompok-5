export interface ApiErrorItem {
  id: string;
  location: string;
  message: string;
  value: unknown;
}

export interface ApiResponse {
  status: number;
  success: boolean;
  title: string;
  detail: string;
  errors: ApiErrorItem[];
}

export function isApiResponse(obj: unknown): obj is ApiResponse {
  if (!obj || typeof obj !== 'object') return false;
  const o = obj as Record<string, unknown>;
  return (
    typeof o.status === 'number' &&
    typeof o.success === 'boolean' &&
    Array.isArray(o.errors)
  );
}
