const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export async function apiRequest(path, options = {}) {
  const token = localStorage.getItem("token");

  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  let response;

  try {
    response = await fetch(`${API_URL}${path}`, {
      ...options,
      headers,
    });
  } catch {
    throw new Error(
      "Tidak dapat terhubung ke server. Periksa apakah backend sedang berjalan."
    );
  }

  const result = await response.json().catch(() => ({}));

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
    }

    throw new Error(
      result.message || `Request gagal (${response.status})`
    );
  }

  return result;
}
