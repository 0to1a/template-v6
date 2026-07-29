import { useState } from 'react';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useMutation } from '@connectrpc/connect-query';

import { requestLogin } from '@/lib/gen/auth/auth-AuthService_connectquery';
import { setPendingEmail } from '@/lib/login';
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

export const Route = createFileRoute('/login/')({
	component: LoginComponent
});

function LoginComponent() {
	const navigate = useNavigate();
	const [email, setEmail] = useState('');
	const [error, setError] = useState('');
	const mutation = useMutation(requestLogin);

	async function submit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault();
		setError('');
		try {
			// Response is generic: never reveals account existence
			await mutation.mutateAsync({ email });
			setPendingEmail(email);
			await navigate({ to: '/login/otp' });
		} catch {
			setError('Something went wrong. Please try again.');
		}
	}

	return (
		<>
			<title>Sign in · Template v6</title>
			<main className="mx-auto flex min-h-screen max-w-sm flex-col justify-center p-8">
				<Card>
					<CardHeader>
						<CardTitle>Sign in</CardTitle>
						<CardDescription>
							Enter your email and we will send you a one-time code.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<form className="flex flex-col gap-4" onSubmit={submit}>
							<div className="flex flex-col gap-1.5">
								<Label htmlFor="email">Email</Label>
								<Input
									id="email"
									type="email"
									name="email"
									required
									value={email}
									onChange={(event) => setEmail(event.target.value)}
								/>
							</div>

							{error && <p className="text-sm text-destructive">{error}</p>}

							<Button type="submit" disabled={mutation.isPending}>
								{mutation.isPending ? 'Sending…' : 'Send code'}
							</Button>
						</form>
					</CardContent>
				</Card>
			</main>
		</>
	);
}
