import { Navigate, Route, Routes } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Payments from "./pages/Payments";
import "./App.css";

function ProtectedRoute({ children }) {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/payments" replace />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route path="/payments"
        element={
          <ProtectedRoute>
            <Payments />
          </ProtectedRoute>
        }
      />

      <Route path="*" element={<Navigate to="/payments" replace />} />
    </Routes>
  );
}
