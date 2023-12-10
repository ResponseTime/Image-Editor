import {react} from "react"
import axios from "axios";
export default function Records(props){
    const handleUpload = async(e)=>{
        const file = e.target.files[0];
        const formData = new FormData();
        formData.append("file", file);
        const res = await axios
        .post("http://localhost:8080/api/v1/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            "Authorization": localStorage.getItem('Auth')
          },
        })
        if(res.data.message){
            console.log("uploaded")
        }
    }
    return <>
    <div className="upload">
        <input name="file" type="file" onChange={handleUpload}/>
    </div>
    <div className="records">

    </div>
    </>
}