import React, { useState } from "react";
import { useApolloClient, useMutation } from "@apollo/client";
import { useRouter } from "next/router";
import { CREATE_USER, LOGIN } from "../lib/graphql/mutations";
import { GET_ME } from "../lib/graphql/queries";

const AuthForm: React.FC = () => {
  const router = useRouter();
  const apolloClient = useApolloClient();
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
          await apolloClient.resetStore();
          await router.push("/dashboard");
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
          await apolloClient.resetStore();
          await router.push("/dashboard");
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
    <div className="card max-w-md mx-auto">
      <h2 className="text-2xl font-semibold mb-6 text-center">
        {isLogin ? "Login" : "Sign Up"}
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {!isLogin && (
          <div>
            <label
              htmlFor="name"
              className="block text-sm font-medium text-slate-600 dark:text-slate-300 mb-1"
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
              className="input"
            />
          </div>
        )}

        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-slate-600 dark:text-slate-300 mb-1"
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
            className="input"
          />
        </div>

        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-slate-600 dark:text-slate-300 mb-1"
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
            className="input"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Processing..." : isLogin ? "Login" : "Sign Up"}
        </button>
      </form>

      {error && (
        <div className="mt-4 rounded-xl border border-rose-200 bg-rose-50 p-3 text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-200">
          Error: {error.message}
        </div>
      )}

      <div className="mt-4 text-center">
        <button
          type="button"
          onClick={() => setIsLogin(!isLogin)}
          className="text-emerald-600 hover:text-emerald-500 text-sm"
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
