import React from "react";
import AuthForm from "../components/AuthForm";

export default function AuthPage() {
  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">
          Authentication
        </h1>
        <p className="text-lg text-gray-600">
          Test user registration and login functionality
        </p>
      </div>

      <AuthForm />

      <div className="bg-green-50 border border-green-200 p-6 rounded-lg max-w-md mx-auto">
        <h3 className="text-lg font-semibold mb-3">🧪 Testing Guide</h3>
        <ul className="list-disc list-inside space-y-2 text-gray-700 text-sm">
          <li>Try creating a new account with the "Sign Up" form</li>
          <li>Test login with existing credentials</li>
          <li>Check how authentication state updates in the header</li>
          <li>JWT tokens are stored in localStorage</li>
          <li>The "me" query will work once authenticated</li>
        </ul>
      </div>
    </div>
  );
}
