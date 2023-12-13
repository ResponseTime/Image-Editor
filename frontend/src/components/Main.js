import React, { useEffect, useState } from "react";
import GridLines from "react-gridlines";
import { motion } from "framer-motion";
import axios from "axios";
export default function Main(props) {
  const handleDownload = async () => {
    const res = await axios.get("http://localhost:8080/api/v1/export", {
      headers: { Authorization: localStorage.getItem("Auth") },
      responseType: "blob",
    });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement("a");
    link.setAttribute("download", res.headers.filename);
    link.href = url;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };
  const [photo, setPhoto] = useState();
  const handleRotate = async () => {
    const res = await axios.post(
      "http://localhost:8080/api/v1/rotate",
      {},
      {
        headers: { Authorization: localStorage.getItem("Auth") },
      }
    );
    setHistory([
      ...history,
      "Image Rotated at" + new Date().toLocaleTimeString(),
    ]);
  };
  const [history, setHistory] = useState([]);
  useEffect(() => {
    const f1 = async () => {
      const res = await axios.get("http://localhost:8080/api/v1/getImage", {
        headers: { Authorization: localStorage.getItem("Auth") },
        responseType: "blob",
      });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      setPhoto(url);
    };
    f1();
  }, [history]);
  return (
    <>
      <div className="editor">
        <GridLines className="editor" cellWidth={400} cellWidth2={200}>
          <motion.img
            style={{
              objectFit: "contain",
              maxHeight: "100%",
              maxWidth: "100%",
            }}
            src={photo}
            alt=""
            drag
            dragTransition={{ bounceStiffness: 600, bounceDamping: 10 }}
            whileTap={{ boxShadow: "0px 0px 15px rgba(0,0,0,0.2)" }}
            dragElastic={0.1}
          />
        </GridLines>
      </div>
      <div className="buttons">
        <button>CROP</button>
        <button>RESIZE</button>
        <button onClick={handleRotate}>ROTATE</button>
        <button>GRAYSCALE</button>
        <button>BLUR</button>
        <button>BRIGHTNESS</button>
        <button>SHARPENING</button>
        <button>contrast</button>
      </div>
      <div className="sidebar">
        <h1 style={{ textAlign: "center" }}>
          Hello {localStorage.getItem("Auth").substring(0, 10)}
        </h1>

        <h1
          style={{
            textAlign: "center",
          }}>
          History
        </h1>
        <div className="history">
          {history.map((Val, index) => {
            return <span key={index}>{Val}</span>;
          })}
        </div>
      </div>
      <div className="exp">
        <button onClick={handleDownload}>Export</button>
        <button>Save</button>
      </div>
    </>
  );
}
