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
    setUtils("rotate");
  };
  const handleCrop = async () => {
    setUtils("crop");
  };
  const [util, setUtils] = useState("");
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
  const rotateLeft = async () => {
    const res = await axios.post(
      "http://localhost:8080/api/v1/rotate",
      {},
      {
        headers: { Authorization: localStorage.getItem("Auth") },
      }
    );
    setHistory([
      ...history,
      "Image Rotated Left at " + new Date().toLocaleTimeString(),
    ]);
  };
  const rotateRight = async () => {
    const res = await axios.post(
      "http://localhost:8080/api/v1/rotater",
      {},
      {
        headers: { Authorization: localStorage.getItem("Auth") },
      }
    );
    setHistory([
      ...history,
      "Image Rotated Right at " + new Date().toLocaleTimeString(),
    ]);
  };
  useEffect(() => {
    setTimeout(() => {
      setUtils("");
    }, 5000);
  }, []);
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
            initial={{ x: 400, y: 400 }}
            drag
            dragTransition={{ bounceStiffness: 600, bounceDamping: 10 }}
            whileTap={{ boxShadow: "0px 0px 15px rgba(0,0,0,0.2)" }}
            dragElastic={0.1}
          />
        </GridLines>
      </div>
      <div className="buttons">
        <button onClick={handleCrop}>CROP</button>
        <button>RESIZE</button>
        <button onClick={handleRotate}>ROTATE</button>
        <button>GRAYSCALE</button>
        <button>BLUR</button>
        <button>BRIGHTNESS</button>
        <button>SHARPENING</button>
        <button>contrast</button>
      </div>
      <div className="sidebar">
        <div className="utildump">
          {util === "crop" ? (
            <div className="crop">
              <button>3:2</button>
              <button>10:9</button>
              <button>16:9</button>
            </div>
          ) : util === "rotate" ? (
            <div className="rot">
              <button onClick={rotateRight}>right</button>
              <button onClick={rotateLeft}>left</button>
            </div>
          ) : (
            <h1></h1>
          )}
        </div>
        <div className="history">
          <h1
            style={{
              textAlign: "center",
            }}>
            History
          </h1>
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
