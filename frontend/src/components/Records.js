import { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import ImageHolder from "./ImageHolder";
export default function Records(props) {
  const navigate = useNavigate();
  const handleUpload = async (e) => {
    const file = e.target.files[0];
    const formData = new FormData();
    formData.append("file", file);
    const res = await axios.post(
      "http://localhost:8080/api/v1/upload",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
          Authorization: localStorage.getItem("Auth"),
        },
      }
    );
    if (res.status === 200) {
      navigate("/main");
    }
  };
  const [data, setData] = useState(null);
  useEffect(() => {
    if (!localStorage.getItem("Auth")) {
      navigate("/login");
      return () => {};
    }
  }, []);
  useEffect(() => {
    const fetchd = async () => {
      try {
        const res = await axios.get("http://localhost:8080/api/v1/getdetails", {
          headers: { Authorization: localStorage.getItem("Auth") },
        });
        const d = await res.data;
        setData(d.data);
      } catch (e) {
        console.log(e);
      }
    };
    fetchd();
  }, []);
  const handleClick = () => {
    navigate("/login");
  };
  const handleClickSign = () => {
    navigate("/signup");
  };
  return (
    <>
      <button onClick={handleClick} className="btn-login">
        Login
      </button>
      <button onClick={handleClickSign} className="btn-login">
        Signup
      </button>
      <div className="upload">
        <input name="file" type="file" onChange={handleUpload} />
      </div>
      <div className="records">
        <h1>Past Projects</h1>
        {data ? (
          data.map((val) => {
            return (
              <div key={val.ProjectName}>
                <ImageHolder pname={val.ProjectName} username={val.User} />
              </div>
            );
          })
        ) : (
          <h1>No Past Projects</h1>
        )}
      </div>
    </>
  );
}
