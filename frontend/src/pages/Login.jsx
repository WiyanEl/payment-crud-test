import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { apiRequest } from "../api/client";

export default function Login() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    email: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  function handleChange(event) {
    setForm({
      ...form,
      [event.target.name]: event.target.value,
    });
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setLoading(true);

    try {
      const result = await apiRequest("/auth/login", {
        method: "POST",
        body: JSON.stringify(form),
      });

      if (!result.data.token) {
        throw new Error("Token tidak ditemukan pada respons server.");
      }

      localStorage.setItem("token", result.data.token);

      const user =
        result.data.user ?? {
          id: result.data.user.id ?? result.data.user.id,
          name: result.data.user.name,
          email: result.data.user.email,
        };

      localStorage.setItem("user", JSON.stringify(user));
      navigate("/payments", { replace: true });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="auth-container">
      <form className="auth-card" onSubmit={handleSubmit}>
        <h1>Payment</h1>
        <p>Masuk untuk mengelola transaksi pembayaran.</p>

        {error && <div className="error-message">{error}</div>}

        <label htmlFor="email">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          autoComplete="email"
          value={form.email}
          onChange={handleChange}
          required
        />

        <label htmlFor="password">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          autoComplete="current-password"
          value={form.password}
          onChange={handleChange}
          required
        />

        <button type="submit" disabled={loading}>
          {loading ? "Memproses..." : "Login"}
        </button>

        <p>
          Belum punya akun? <Link to="/register">Daftar di sini</Link>
        </p>
      </form>
    </main>
  );
}
