import { useEffect, useState } from 'react';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useMutation } from '@connectrpc/connect-query';

import { submitLogin } from '@/lib/gen/auth/auth-AuthService_connectquery';
import { completeLogin, getPendingEmail } from '@/lib/login';
import { Button } from '@/lib/components/ui/button';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle
} from '@/lib/components/ui/card';
import { Input } from '@/lib/components/ui/input';
import { Label } from '@/lib/components/ui/label';

export const Route = createFileRoute('/login/otp')({
	component: OtpComponent
});

function OtpComponent() {
	const navigate = useNavigate();
	const email = getPendingEmail();

	// No email means hard refresh; restart the flow
	useEffect(() => {
		if (!email) void navigate({ to: '/login' });
	}, [email, navigate]);

	const [code, setCode] = useState('');
	const [error, setError] = useState('');
	const mutation = useMutation(submitLogin);

	async function submit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault();
		setError('');
		try {
			await completeLogin(
				email,
				code,
				(path) => navigate({ to: path as '/' }),
				mutation.mutateAsync
			);
		} catch {
			// Invalid email and code deliberately look the same
			setError('Invalid email or code.');
		}
	}

	return (
		<>
			<title>Enter your code · Template v6</title>
			<main className="mx-auto flex min-h-screen max-w-sm flex-col justify-center p-8">
				<Card>
					<CardHeader>
						<CardTitle>Enter your code</CardTitle>
						<CardDescription>
							We sent a one-time code to {email}. It is valid for 5 minutes.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<form className="flex flex-col gap-4" onSubmit={submit}>
							<div className="flex flex-col gap-1.5">
								<Label htmlFor="code">One-time code</Label>
								<Input
									id="code"
									className="tracking-widest"
									type="text"
									name="code"
									inputMode="numeric"
									autoComplete="one-time-code"
									required
									value={code}
									onChange={(event) => setCode(event.target.value)}
								/>
							</div>

							{error && <p className="text-sm text-destructive">{error}</p>}

							<Button type="submit" disabled={mutation.isPending}>
								{mutation.isPending ? 'Signing in…' : 'Sign in'}
							</Button>
						</form>
					</CardContent>
				</Card>
			</main>
		</>
	);
}
