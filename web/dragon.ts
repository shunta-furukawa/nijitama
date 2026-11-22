export function dragon(c: string, grown = false) {
  return `<svg role="img" aria-label="${
    grown ? "つばさが育ったドラゴン" : "小さなドラゴン"
  }" viewBox="0 0 180 124"><ellipse cx="91" cy="111" rx="55" ry="7" fill="#233d3e" opacity=".12"/><g fill="${c}" stroke="#354e50" stroke-width="2.5" stroke-linejoin="round"><path d="M86 91Q49 108 22 72Q20 108 62 108L92 106Z"/><path d="M78 82Q${
    grown ? "20 2 31 49" : "39 28 43 58"
  } 44 73Q57 67 60 84Z"/><ellipse cx="87" cy="86" rx="33" ry="25"/><path d="M88 47L83 22L101 41M117 41L132 20L129 48" fill="#efd99b"/><path d="M81 60Q80 33 106 34Q128 31 135 53Q158 52 154 68Q147 84 112 81Q86 84 81 60Z"/><path d="M68 97Q58 115 76 111L82 102M96 99Q91 115 108 110L111 99"/></g><ellipse cx="118" cy="52" rx="4" ry="6" fill="#273f43"/><circle cx="119" cy="50" r="1.5" fill="white"/><circle cx="145" cy="59" r="2" fill="#354e50"/><path d="M133 70q6 4 12-1" fill="none" stroke="#354e50" stroke-width="2"/><ellipse cx="115" cy="65" rx="7" ry="3.5" fill="#ecac92"/></svg>`;
}
