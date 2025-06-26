import React, { useState } from "react";
import { useMutation } from "@apollo/client";
import { CREATE_USER, LOGIN } from "../lib/graphql/mutations";
import { GET_ME } from "../lib/graphql/queries";

const AuthForm: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
  });

  const [createUser, { loading: createLoading, error: createError }] =
    useMutation(CREATE_USER);
  const [login, { loading: loginLoading, error: loginError }] = useMutation(
    LOGIN,
    {
      refetchQueries: [{ query: GET_ME }],
    }
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (isLogin) {
        const { data } = await login({
          variables: {
            email: formData.email,
            password: formData.password,
          },
        });

        if (data?.login) {
          localStorage.setItem("auth-token", data.login);
          window.location.reload();
        }
      } else {
        await createUser({
          variables: {
            name: formData.name,
            email: formData.email,
            password: formData.password,
          },
        });

        // Auto-login after registration
        const { data } = await login({
          variables: {
            email: formData.email,
            password: formData.password,
          },
        });

        if (data?.login) {
          localStorage.setItem("auth-token", data.login);
          window.location.reload();
        }
      }
    } catch (err) {
      console.error("Auth error:", err);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const loading = createLoading || loginLoading;
  const error = createError || loginError;

  return (
    <div className="bg-white p-6 rounded-lg shadow-md max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-6 text-center">
        {isLogin ? "Login" : "Sign Up"}
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {!isLogin && (
          <div>
            <label
              htmlFor="name"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Name
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              required={!isLogin}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        )}

        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Email
          </label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleInputChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Password
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleInputChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Processing..." : isLogin ? "Login" : "Sign Up"}
        </button>
      </form>

      {error && (
        <div className="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          Error: {error.message}
        </div>
      )}

      <div className="mt-4 text-center">
        <button
          type="button"
          onClick={() => setIsLogin(!isLogin)}
          className="text-blue-600 hover:text-blue-800 text-sm"
        >
          {isLogin
            ? "Don't have an account? Sign up"
            : "Already have an account? Login"}
        </button>
      </div>
    </div>
  );
};

export default AuthForm;
