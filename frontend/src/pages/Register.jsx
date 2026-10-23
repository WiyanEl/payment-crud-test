import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { apiRequest } from "../api/client";

export default function Register() {
  const navigate = useNavigate();

  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  function handleChange(event) {
    setForm({
      ...form,
      [event.target.name]: event.target.value,
    });
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSuccess("");

    if (form.password.length < 8) {
      setError("Password minimal 8 karakter.");
      return;
    }

    setLoading(true);

    try {
      await apiRequest("/auth/register", {
        method: "POST",
        body: JSON.stringify(form),
      });

      setSuccess("Registrasi berhasil. Mengarahkan ke halaman login...");
      navigate("/login", {
        replace: true,
        state: { registered: true },
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="auth-container">
      <form className="auth-card" onSubmit={handleSubmit}>
        <h1>Buat Akun</h1>
        <p>Daftar untuk mulai mengelola transaksi.</p>

        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}

        <label htmlFor="name">Nama</label>
        <input
          id="name"
          name="name"
          type="text"
          autoComplete="name"
          value={form.name}
          onChange={handleChange}
          required
        />

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
          autoComplete="new-password"
          minLength={8}
          value={form.password}
          onChange={handleChange}
          required
        />

        <button type="submit" disabled={loading}>
          {loading ? "Mendaftarkan..." : "Daftar"}
        </button>

        <p>
          Sudah punya akun? <Link to="/login">Login di sini</Link>
        </p>
      </form>
    </main>
  );
}
