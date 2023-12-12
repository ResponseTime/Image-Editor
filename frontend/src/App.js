import { useEffect, useState } from "react";
import "./App.css";
import Auth from "./components/Auth";
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  useNavigate,
} from "react-router-dom";
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
        <Route
          path="/"
          element={auth ? <Records /> : <Navigate to="/login" />}
        />
        <Route path="/main" element={<Main />} />
        {/* <Route path="/upload" element={<Records />} /> */}
        <Route path="/login" element={<Auth />} />
        <Route path="/signup" element={<Auth title="signup" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
