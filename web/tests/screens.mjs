// The Go test bridges fetch requests to the real HTTP handler and a fresh SQLite database.
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import vm from "node:vm";
import { createInterface } from "node:readline";
const { parseHTML } = await import(process.env.NIJI_DOM_MODULE || "linkedom");
const { window } = parseHTML(
  '<html><body><main id="app"></main></body></html>'
);
const document = window.document;
const catalog = JSON.parse(readFileSync("../curriculum/islands.json"));
const lines = createInterface({ input: process.stdin })[Symbol.asyncIterator]();
async function fetch(url, init) {
  if (url === "/catalog") return { ok: true, json: async () => catalog };
  process.stdout.write(init.body + "\n");
  const { status, body } = JSON.parse((await lines.next()).value);
  return {
    ok: status === 200,
    json: async () => JSON.parse(body),
    text: async () => body,
  };
}
let fields = {};
const context = vm.createContext({
  window,
  document,
  navigator: { userAgent: "DOM test" },
  localStorage: {},
  fetch,
  console,
  setTimeout,
  clearTimeout,
  FormData: class {
    get(name) {
      return fields[name];
    }
  },
});
context.self = context;
window.scrollTo = () => {};
for (const f of ["vendor/react.js", "vendor/react-dom.js", "app.js"])
  vm.runInContext(readFileSync("standalone/" + f, "utf8"), context);
const tick = () => new Promise((r) => setTimeout(r, 120));
async function click(text) {
  const b = [...document.querySelectorAll("button")].find((b) =>
    b.textContent.includes(text)
  );
  assert.ok(b, text + " " + document.body.textContent);
  b.click();
  await tick();
}
async function submit() {
  document
    .querySelector("form")
    .dispatchEvent(
      new window.Event("submit", { bubbles: true, cancelable: true })
    );
  await tick();
}
await tick();
await click("保護者の部屋");
fields = { pin: process.env.NIJI_TEST_PIN };
await submit();
fields = { name: "ソラ", pin: "5678", grade: "2", icon: "0" };
await submit();
await click("人をえらぶ");
await click("ソラ");
await submit();
assert.equal(document.querySelectorAll(".card").length, 5);
await click("小数のかけ橋へ");
await click("ヒント");
assert.match(document.body.textContent, /小数第２位/);
await click("橋をわたる");
assert.match(document.body.textContent, /石を足して/);
for (let round = 0; round < 5; round++) {
  const target = document
    .querySelector(".game p strong")
    .textContent.replace(".", "")
    .split("")
    .map(Number);
  for (const [i, place] of ["1", "0.1", "0.01"].entries()) {
    for (let j = 0; j < target[i]; j++) {
      document
        .querySelector('[aria-label="' + place + 'の石をふやす"]')
        .click();
      await new Promise((r) => setTimeout(r, 5));
    }
  }
  await click("橋をわたる");
  if (round < 4) await click("つぎの橋");
}
assert.match(document.body.textContent, /谷をわたりきった/);
await click("島にもどる");
assert.match(document.body.textContent, /3 エレメント/);
await click("人をえらぶ");
await click("保護者の部屋");
fields = { pin: process.env.NIJI_TEST_PIN };
await submit();
assert.match(document.body.textContent, /正解 5\/6回/);
process.exit(0);
