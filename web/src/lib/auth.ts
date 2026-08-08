// Only module for the token; localStorage has XSS risk
const ACCESS_TOKEN_KEY = 'template-v6.access_token';

// In-memory fallback for non-browser environments (tests)
let memoryToken: string | null = null;

function hasLocalStorage(): boolean {
	return typeof localStorage !== 'undefined';
}

export function getAccessToken(): string | null {
	return hasLocalStorage() ? localStorage.getItem(ACCESS_TOKEN_KEY) : memoryToken;
}

export function setAccessToken(token: string): void {
	if (hasLocalStorage()) localStorage.setItem(ACCESS_TOKEN_KEY, token);
	else memoryToken = token;
}

export function clearAccessToken(): void {
	if (hasLocalStorage()) localStorage.removeItem(ACCESS_TOKEN_KEY);
	else memoryToken = null;
}

export function isAuthenticated(): boolean {
	return getAccessToken() !== null;
}

// No server session; JWT stays valid but unused after this
export function logout(): void {
	clearAccessToken();
}
