import { useEffect, useState } from "react";

const initialForm = {
  customer_number: "",
  customer_name: "",
  amount: "",
  method: "transfer",
  status: "pending",
};

export default function PaymentForm({
  initialData,
  onSubmit,
  onCancel,
  loading,
}) {
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState("");

  useEffect(() => {
    if (initialData) {
      setForm({
        customer_number: initialData.customer_number ?? "",
        customer_name: initialData.customer_name ?? "",
        amount: String(initialData.amount ?? ""),
        method: initialData.method ?? "transfer",
        status: initialData.status ?? "pending",
      });
    } else {
      setForm(initialForm);
    }

    setError("");
  }, [initialData]);

  function handleChange(event) {
    setForm((current) => ({
      ...current,
      [event.target.name]: event.target.value,
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");

    const amount = Number(form.amount);

    if (
      !form.customer_number.trim() ||
      !form.customer_name.trim() ||
      !Number.isSafeInteger(amount) ||
      amount <= 0
    ) {
      setError("Lengkapi data dan masukkan nominal positif berupa bilangan bulat.");
      return;
    }

    try {
      await onSubmit({
        customer_number: form.customer_number.trim(),
        customer_name: form.customer_name.trim(),
        amount,
        method: form.method,
        status: form.status,
      });
    } catch (err) {
      setError(err.message || "Gagal menyimpan transaksi.");
    }
  }

  return (
    <form className="payment-form" onSubmit={handleSubmit}>
      <h2>{initialData ? "Edit Transaksi" : "Tambah Transaksi"}</h2>

      {error && <div className="error-message">{error}</div>}

      <label htmlFor="customer_number">Nomor Pelanggan</label>
      <input
        id="customer_number"
        name="customer_number"
        value={form.customer_number}
        onChange={handleChange}
        required
      />

      <label htmlFor="customer_name">Nama Pelanggan</label>
      <input
        id="customer_name"
        name="customer_name"
        value={form.customer_name}
        onChange={handleChange}
        required
      />

      <label htmlFor="amount">Nominal (Rp)</label>
      <input
        id="amount"
        name="amount"
        type="number"
        min="1"
        step="1"
        value={form.amount}
        onChange={handleChange}
        required
      />

      <label htmlFor="method">Metode Pembayaran</label>
      <select
        id="method"
        name="method"
        value={form.method}
        onChange={handleChange}
      >
        <option value="transfer">Transfer</option>
        <option value="ewallet">E-Wallet</option>
        <option value="va">Virtual Account</option>
      </select>

      {initialData && (
        <>
          <label htmlFor="status">Status</label>
          <select
            id="status"
            name="status"
            value={form.status}
            onChange={handleChange}
          >
            <option value="pending">Pending</option>
            <option value="success">Success</option>
            <option value="failed">Failed</option>
          </select>
        </>
      )}

      <div className="form-actions">
        <button type="submit" disabled={loading}>
          {loading ? "Menyimpan..." : "Simpan"}
        </button>
        <button
          type="button"
          className="button-secondary"
          onClick={onCancel}
          disabled={loading}
        >
          Batal
        </button>
      </div>
    </form>
  );
}