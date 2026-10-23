import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiRequest } from "../api/client";
import PaymentForm from "../components/PaymentForm";

const PAGE_SIZE = 10;

function formatRupiah(amount) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(value) {
  if (!value) return "-";

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) return "-";

  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

export default function Payments() {
  const navigate = useNavigate();

  const [payments, setPayments] = useState([]);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: PAGE_SIZE,
    total: 0,
    total_pages: 0,
  });

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [formOpen, setFormOpen] = useState(false);
  const [editingPayment, setEditingPayment] = useState(null);

  const user = (() => {
    try {
      return JSON.parse(localStorage.getItem("user") || "null");
    } catch {
      return null;
    }
  })();

  const loadPayments = useCallback(async (page = 1) => {
    setLoading(true);
    setError("");

    try {
      const result = await apiRequest(
        `/payments?page=${page}&limit=${PAGE_SIZE}`
      );

      setPayments(result.data ?? []);
      setPagination(result.pagination ?? {
        page,
        limit: PAGE_SIZE,
        total: 0,
        total_pages: 0,
      });
    } catch (err) {
      if (!localStorage.getItem("token")) {
        navigate("/login", { replace: true });
        return;
      }

      setError(err.message || "Gagal memuat transaksi.");
    } finally {
      setLoading(false);
    }
  }, [navigate]);

  useEffect(() => {
    loadPayments(1);
  }, [loadPayments]);

  function openCreateForm() {
    setEditingPayment(null);
    setNotice("");
    setFormOpen(true);
  }

  function openEditForm(payment) {
    setEditingPayment(payment);
    setNotice("");
    setFormOpen(true);
  }

  async function handleSave(formData) {
    setSaving(true);
    setNotice("");

    try {
      if (editingPayment) {
        await apiRequest(`/payments/${editingPayment.id}`, {
          method: "PUT",
          body: JSON.stringify(formData),
        });
        setNotice("Transaksi berhasil diperbarui.");
      } else {
        await apiRequest("/payments", {
          method: "POST",
          body: JSON.stringify({
            customer_number: formData.customer_number,
            customer_name: formData.customer_name,
            amount: formData.amount,
            method: formData.method,
          }),
        });
        setNotice("Transaksi berhasil ditambahkan.");
      }

      setFormOpen(false);
      setEditingPayment(null);

      // Setelah create, kembali ke halaman pertama agar data terbaru terlihat.
      await loadPayments(1);
    } catch (err) {
      if (!localStorage.getItem("token")) {
        navigate("/login", { replace: true });
      }
      throw err;
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(payment) {
    const confirmed = window.confirm(
      `Hapus transaksi ${payment.customer_number} - ${payment.customer_name}?`
    );

    if (!confirmed) return;

    setError("");
    setNotice("");

    try {
      await apiRequest(`/payments/${payment.id}`, {
        method: "DELETE",
      });

      setNotice("Transaksi berhasil dihapus.");

      const nextPage =
        payments.length === 1 && pagination.page > 1
          ? pagination.page - 1
          : pagination.page;

      await loadPayments(nextPage);
    } catch (err) {
      if (!localStorage.getItem("token")) {
        navigate("/login", { replace: true });
        return;
      }

      setError(err.message || "Gagal menghapus transaksi.");
    }
  }

  function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    navigate("/login", { replace: true });
  }

  return (
    <main className="page-container">
      <header className="dashboard-header">
        <div>
          <h1>Payment Transactions</h1>
          <p>
            Selamat datang{user?.name ? `, ${user.name}` : ""}.
            Kelola transaksi pembayaranmu di sini.
          </p>
        </div>

        <button
          type="button"
          className="button-secondary"
          onClick={handleLogout}
        >
          Logout
        </button>
      </header>

      {notice && <div className="success-message">{notice}</div>}
      {error && (
        <div className="error-message">
          {error}
          <button
            type="button"
            className="text-button"
            onClick={() => loadPayments(pagination.page)}
          >
            Coba lagi
          </button>
        </div>
      )}

      {!formOpen && (
        <div className="toolbar">
          <div>
            <strong>{pagination.total}</strong> transaksi
          </div>
          <button type="button" onClick={openCreateForm}>
            + Tambah Transaksi
          </button>
        </div>
      )}

      {formOpen && (
        <section className="form-panel">
          <PaymentForm
            key={editingPayment?.id ?? "new-payment"}
            initialData={editingPayment}
            onSubmit={handleSave}
            onCancel={() => {
              setFormOpen(false);
              setEditingPayment(null);
            }}
            loading={saving}
          />
        </section>
      )}

      <section className="table-panel">
        {loading ? (
          <div className="state-message">Memuat transaksi...</div>
        ) : payments.length === 0 ? (
          <div className="state-message">
            <h2>Belum ada transaksi</h2>
            <p>Tambahkan transaksi pertamamu untuk mulai mengelola pembayaran.</p>
            {!formOpen && (
              <button type="button" onClick={openCreateForm}>
                Tambah Transaksi
              </button>
            )}
          </div>
        ) : (
          <>
            <div className="table-scroll">
              <table className="payments-table">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Pelanggan</th>
                    <th>Nominal</th>
                    <th>Metode</th>
                    <th>Status</th>
                    <th>Dibuat</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {payments.map((payment, index) => (
                    <tr key={payment.id}>
                      <td>{index + 1}</td>
                      <td>
                        <strong>{payment.customer_name}</strong>
                        <span className="secondary-text">
                          {payment.customer_number}
                        </span>
                      </td>
                      <td>{formatRupiah(payment.amount)}</td>
                      <td>{payment.method}</td>
                      <td>
                        <span
                          className={`status-badge status-${payment.status}`}
                        >
                          {payment.status}
                        </span>
                      </td>
                      <td>{formatDate(payment.created_at)}</td>
                      <td>
                        <div className="action-buttons">
                          <button
                            type="button"
                            className="button-small button-secondary"
                            onClick={() => openEditForm(payment)}
                          >
                            Edit
                          </button>
                          <button
                            type="button"
                            className="button-small button-danger"
                            onClick={() => handleDelete(payment)}
                          >
                            Hapus
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <footer className="pagination">
              <span>
                Halaman {pagination.page} dari{" "}
                {Math.max(pagination.total_pages, 1)}
              </span>

              <div className="action-buttons">
                <button
                  type="button"
                  className="button-secondary"
                  disabled={pagination.page <= 1 || loading}
                  onClick={() => loadPayments(pagination.page - 1)}
                >
                  Sebelumnya
                </button>
                <button
                  type="button"
                  className="button-secondary"
                  disabled={
                    pagination.page >= pagination.total_pages || loading
                  }
                  onClick={() => loadPayments(pagination.page + 1)}
                >
                  Berikutnya
                </button>
              </div>
            </footer>
          </>
        )}
      </section>
    </main>
  );
}
