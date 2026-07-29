import { useState } from 'react';
import { createFileRoute, Link } from '@tanstack/react-router';

import { isAuthenticated, logout } from '@/lib/auth';
import { Button } from '@/lib/components/ui/button';

export const Route = createFileRoute('/')({
	component: HomeComponent
});

function HomeComponent() {
	const [authed, setAuthed] = useState(isAuthenticated);

	function handleLogout() {
		logout();
		setAuthed(false);
	}

	return (
		<>
			<title>Template v6</title>
			<main className="mx-auto flex min-h-screen max-w-md flex-col items-center justify-center gap-4 p-8">
				<h1 className="text-2xl font-semibold">Template v6</h1>

				{authed ? (
					<>
						<p className="text-muted-foreground">You are signed in.</p>
						<Button type="button" onClick={handleLogout}>
							Log out
						</Button>
					</>
				) : (
					<>
						<p className="text-muted-foreground">You are not signed in.</p>
						<Button render={<Link to="/login" />}>Sign in</Button>
					</>
				)}
			</main>
		</>
	);
}
