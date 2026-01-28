import React from "react";
import AuthForm from "../components/AuthForm";

export default function AuthPage() {
  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-semibold mb-4">Authentication</h1>
        <p className="text-lg text-slate-600 dark:text-slate-300">
          Login or create an account to personalize Muse.
        </p>
      </div>

      <AuthForm />

      <div className="card-muted max-w-md mx-auto">
        <h3 className="text-lg font-semibold mb-3">Testing guide</h3>
        <ul className="list-disc list-inside space-y-2 text-slate-600 dark:text-slate-300 text-sm">
          <li>Try creating a new account with the "Sign Up" form</li>
          <li>Test login with existing credentials</li>
          <li>Check how authentication state updates in the header</li>
          <li>Successful login redirects you to the dashboard</li>
          <li>JWT tokens are stored in localStorage</li>
          <li>The "me" query will work once authenticated</li>
        </ul>
      </div>
    </div>
  );
}
