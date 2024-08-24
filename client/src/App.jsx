import React from "react";
import "./App.css";

function App() {
  return (
    <div className="container">
      <h1>Aplicación de Tareas</h1>
      <button onClick={async()=> {
        const backendUrl = process.env.REACT_APP_BACKEND_URL || '';
        const res = await fetch(`${backendUrl}/works`);        
        const data = await res.json()
        console.log(data)
      }} >Obtener Tareas</button>
    </div>
  );
}

export default App;
