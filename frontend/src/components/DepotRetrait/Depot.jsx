import React from "react";
import "./depot.css";
const Depot = () => {
  return (
    <div className="depotRetrait">
      <h4> le lieux ou vous etes entraint de faire un depot</h4>
      <input
        type="text"
        placeholder="entrer le lieux ou vous etes entraint de faire "
      />
      <h4>valeur de que vous etes entraint de faire un depot</h4>
      <input
        type="number"
        name="valueDep"
        id="valueDep"
        placeholder="la valeur de que vous etes entraint de faire un depot"
      />
      <h4>votre mots de passes</h4>
      <input
        type="password"
        name="numbPass"
        id="numbPas"
        placeholder="votre mots de passes"
      />
      <button> confirmer </button>
    </div>
  );
};

export default Depot;