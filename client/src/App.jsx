import React from "react";
import "./App.css";

function App() {
  return (
    <div className="container">
      <h1>Aplicación de Tareas</h1>
      <button onClick={async()=> {
        const res = await fetch('http://localhost:3030/works')
        const data = await res.json()
        console.log(data)
      }} >Obtener Tareas</button>
    </div>
  );
}

export default App;
