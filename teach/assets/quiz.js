// quiz.js — composant de quiz réutilisable pour toutes les leçons.
// Boucle de feedback immédiate (récupération en mémoire = mémoire long terme).
//
// Usage dans une leçon :
//   <div class="quiz" data-answer="1">
//     <div class="q">La question ?</div>
//     <button>Mauvaise réponse</button>
//     <button>Bonne réponse</button>   <!-- index 1 (0-based) = data-answer -->
//     <button>Mauvaise réponse</button>
//     <div class="explain">Explication montrée après la réponse.</div>
//   </div>
// Convention pédagogique : les réponses ont la MÊME longueur (aucun indice de forme).

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".quiz").forEach((quiz) => {
    const answer = parseInt(quiz.dataset.answer, 10);
    const buttons = Array.from(quiz.querySelectorAll("button"));
    const explain = quiz.querySelector(".explain");
    let done = false;

    buttons.forEach((btn, i) => {
      btn.addEventListener("click", () => {
        if (done) return;
        done = true;
        buttons[answer]?.classList.add("correct");
        if (i !== answer) btn.classList.add("wrong");
        explain?.classList.add("show");
        buttons.forEach((b) => (b.style.cursor = "default"));
      });
    });
  });
});
