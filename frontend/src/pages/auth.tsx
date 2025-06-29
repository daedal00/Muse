import React from "react";
import AuthForm from "../components/AuthForm";
import Layout from "../components/Layout/Layout";

export default function AuthPage() {
  return (
    <Layout>
      <div className="max-w-md mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            Welcome to Muse
          </h1>
          <p className="text-gray-600">Sign in or create an account</p>
        </div>
        <AuthForm />
      </div>
    </Layout>
  );
}
