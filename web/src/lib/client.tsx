// Only place that builds transport + bearer header
import type { ReactNode } from 'react';
import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { TransportProvider } from '@connectrpc/connect-query';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { clearAccessToken, getAccessToken } from './auth';

const authInterceptor: Interceptor = (next) => async (req) => {
	const token = getAccessToken();
	if (token) {
		req.header.set('Authorization', `Bearer ${token}`);
	}

	try {
		return await next(req);
	} catch (err) {
		const connectError = ConnectError.from(err);
		if (connectError.code === Code.Unauthenticated) {
			clearAccessToken();
		}
		throw connectError;
	}
};

export const transport = createConnectTransport({
	baseUrl: '/',
	interceptors: [authInterceptor]
});

export const queryClient = new QueryClient();

// Mount AppProviders once at the root, above the router.
export function AppProviders({ children }: { children: ReactNode }) {
	return (
		<TransportProvider transport={transport}>
			<QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
		</TransportProvider>
	);
}
