import { useEffect, useState } from "react";
import "./App.css";
import Auth from "./components/Auth";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Main from "./components/Main";
import Records from "./components/Records";
function App() {
  return (
	<BrowserRouter>
      <Routes>
        <Route path="/" element={<Records />} />
        <Route path="/main" element={<Main />} />
        {/* <Route path="/upload" element={<Records />} /> */}
        <Route path="/login" element={<Auth />} />
        <Route path="/signup" element={<Auth title="signup" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
