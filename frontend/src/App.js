import { useEffect, useState } from "react";
import "./App.css";
import Auth from "./components/Auth";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Main from "./components/Main";
import Records from "./components/Records";
function App() {
  const [auth, setAuth] = useState(true);
  useEffect(() => {
    if (!localStorage.getItem("Auth")) {
      setAuth(false);
    } else {
      setAuth(true);
    }
  }, [auth]);
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={auth ? <Main /> : <Navigate to="/login" />} />
        <Route path="/upload" element={<Records />} />
        <Route path="/login" element={<Auth />} />
        <Route path="/signup" element={<Auth title="signup" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
