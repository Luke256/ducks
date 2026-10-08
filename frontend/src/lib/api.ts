export const apiBase = __API_URL__.replace(/\/+$/, "");
export async function api<T = void>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  let response: Response;
  try {
    response = await fetch(`${apiBase}${path}`, options);
  } catch (error) {
    if (
      error &&
      typeof error === "object" &&
      "name" in error &&
      error.name === "AbortError"
    )
      throw error;
    throw new Error("通信できませんでした。接続を確認してください。");
  }
  if (!response.ok) {
    const detail = await response.text();
    let message = detail;
    try {
      message = JSON.parse(detail).message || detail;
    } catch {
      /* テキスト形式のAPIエラーもある。 */
    }
    throw new Error(
      `処理に失敗しました (${response.status})${message ? `: ${message}` : ""}`,
    );
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}
export const jsonBody = (body: unknown): RequestInit => ({
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(body),
});
export function imageUrl(path: string): string {
  if (!path) return "";
  return new URL(path, new URL(`${apiBase}/`, window.location.origin)).href;
}
