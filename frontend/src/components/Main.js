import React, { useEffect, useState } from "react";
import GridLines from "react-gridlines";
import { motion } from "framer-motion";
import axios from "axios";
import { useNavigate } from "react-router-dom";
export default function Main(props) {
  const navigate = useNavigate();
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
  const [saveText, setSaveText] = useState("");
  const saveCall = () => {
    setUtils("save");
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
  const handleExit = () => {
    localStorage.removeItem("Auth");
    navigate("/");
  };
  useEffect(() => {
    setInterval(() => {
      setUtils("");
    }, 50000);
  }, []);
  const handleBlur = async () => {
    const res = await axios.get("http://localhost:8080/api/v1/blurinc", {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      "Image Blured At " + new Date().toLocaleTimeString(),
    ]);
  };
  const handleSharp = async () => {
    const res = await axios.get("http://localhost:8080/api/v1/sharpinc", {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      "Image Sharpend at " + new Date().toLocaleTimeString(),
    ]);
  };
  const handleSave = async () => {
    const res = await axios.get(
      `http://localhost:8080/api/v1/save/${saveText}`,
      {
        headers: { Authorization: localStorage.getItem("Auth") },
      }
    );
    setHistory([
      ...history,
      `${saveText} Saved at ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const handleResize = () => {
    setUtils("resize");
  };
  const handleGrayscale = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/grayscale`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Grayscaled at ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const handleBrightUp = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/brightinc`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Brigtness Increased At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const handleBrightDown = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/brightdec`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Brigtness Decreased At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const handleContrastUp = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/contrastinc`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Contrast Increased At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const handleContrastDown = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/contrastdec`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Contrast Decreased At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const re1 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/re1`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Resized At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const re2 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/re2`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Resized At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const re3 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/re3`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Resized At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const c1 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/c1`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Cropped At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const c2 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/c2`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Cropped At ${new Date().toLocaleTimeString()}`,
    ]);
  };
  const c3 = async () => {
    const res = await axios.get(`http://localhost:8080/api/v1/c3`, {
      headers: { Authorization: localStorage.getItem("Auth") },
    });
    setHistory([
      ...history,
      `Image Cropped At ${new Date().toLocaleTimeString()}`,
    ]);
  };
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
        <button onClick={handleResize}>RESIZE</button>
        <button onClick={handleRotate}>ROTATE</button>
        <button onClick={handleGrayscale}>GRAYSCALE</button>
        <button onClick={handleBlur}>BLUR</button>
        <button onClick={handleSharp}>SHARPENING</button>
        <button
          onClick={() => {
            setUtils("bright");
          }}>
          BRIGHTNESS
        </button>
        <button
          onClick={() => {
            setUtils("contrast");
          }}>
          contrast
        </button>
      </div>
      <div className="sidebar">
        <div className="utildump">
          {util === "crop" && (
            <div className="crop">
              <button onClick={c1}>3:2</button>
              <button onClick={c2}>10:9</button>
              <button onClick={c3}>16:9</button>
            </div>
          )}

          {util === "rotate" && (
            <div className="rot">
              <button onClick={rotateRight}>right</button>
              <button onClick={rotateLeft}>left</button>
            </div>
          )}

          {util === "" && (
            <>
              <h1 style={{ textAlign: "center" }}>ImageCraft</h1>
              <img
                style={{ objectFit: "cover", height: "200px", width: "100%" }}
                src="https://static-gcp.freepikcompany.com/web-app/media/wepik-1-2000.webp"
                alt=""
              />
            </>
          )}
          {util === "save" && (
            <div className="save">
              <input
                type="text"
                value={saveText}
                onChange={(e) => {
                  setSaveText(e.target.value);
                }}
              />
              <button onClick={handleSave}>Save Project</button>
            </div>
          )}
          {util === "resize" && (
            <div className="crop">
              <button onClick={re1}>800X400</button>
              <button onClick={re2}>1920X1080</button>
              <button onClick={re3}>1280x1024</button>
            </div>
          )}
          {util === "bright" && (
            <div className="crop">
              <button onClick={handleBrightUp}>Increase Brightness</button>
              <button onClick={handleBrightDown}>Decrease Brightness</button>
            </div>
          )}
          {util === "contrast" && (
            <div className="crop">
              <button onClick={handleContrastUp}>Increase Contrast</button>
              <button onClick={handleContrastDown}>Decrease Contrast</button>
            </div>
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
      <button className="exit" onClick={handleExit}>
        Exit and Logout
      </button>
      <div className="exp">
        <button onClick={handleDownload}>Export</button>
        <button onClick={saveCall}>Save</button>
      </div>
    </>
  );
}
