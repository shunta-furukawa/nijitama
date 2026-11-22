import * as React from "react";
import { createRoot } from "react-dom/client";
import { dragon } from "./dragon";
type Learner = { id: string; name: string; grade: number; icon: number };
type Game = {
  run: string;
  step: number;
  target: number;
  hint: number;
  correct?: boolean;
  reward?: number;
};
type RecordRow = {
  owner: string;
  name: string;
  started: number;
  ended: number | null;
  run: string | null;
  step: number;
  answers: number;
  correct: number;
  seconds: number;
  mastery: number;
};
type Catalog = {
  worlds: [string, string, string][];
  unit: { rounds: number; name: string };
};
type Page =
  | "profiles"
  | "parent"
  | "select"
  | "report"
  | "home"
  | "game"
  | "result";
const icons = ["🌱", "🌙", "🍄", "🌻", "🐚"];
const decimal = (n: number) => (n / 100).toFixed(2);
async function api<T>(action: string, fields: object = {}): Promise<T> {
  const response = await fetch("/api", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ action, ...fields }),
  });
  if (!response.ok) throw Error(await response.text());
  return response.json();
}
function Dragon({
  color = "#84c9d1",
  grown = false,
}: {
  color?: string;
  grown?: boolean;
}) {
  return (
    <div
      className="dragon"
      dangerouslySetInnerHTML={{ __html: dragon(color, grown) }}
    />
  );
}
function Hero({
  title,
  caption,
  children,
}: {
  title: string;
  caption: string;
  children: React.ReactNode;
}) {
  return (
    <section className="hero">
      <div>
        <small>{caption}</small>
        <h1>{title}</h1>
        {children}
      </div>
      <div className="scene" aria-hidden="true">
        <svg
          className="landscape"
          viewBox="0 0 380 190"
          preserveAspectRatio="xMidYMax slice"
        >
          <defs>
            <pattern
              id="grass"
              width="24"
              height="20"
              patternUnits="userSpaceOnUse"
            >
              <path d="m4 12 2-4 1 4 3-3" fill="none" stroke="#8eab82" />
            </pattern>
          </defs>
          <circle cx="280" cy="42" r="25" fill="#f5d795" />
          <g fill="#fbfaf0">
            <path d="M12 48q9-15 20-6 8-17 24-3 16-5 24 9zM288 78q9-15 20-6 8-17 24-3 16-5 24 9z" />
          </g>
          <path
            d="M0 130Q60 45 130 112Q221 26 314 105Q352 65 380 87V190H0"
            fill="#b3cbae"
          />
          <path d="M0 132Q90 101 170 134Q264 92 380 119V190H0" fill="#9fbc95" />
          <path
            d="M0 132Q90 101 170 134Q264 92 380 119V190H0"
            fill="url(#grass)"
          />
          <path
            d="M223 113Q136 129 187 137T121 190H205Q262 159 215 147T223 113"
            fill="#bce3dd"
          />
          {[28, 58, 304, 332, 358].map((x, i) => (
            <g
              key={x}
              transform={"translate(" + x + " " + (98 + (i % 2) * 9) + ")"}
            >
              <path d="M0 35V8" stroke="#776b51" strokeWidth="4" />
              <path
                d="m-14 23 8-13-5 2L0-10 12 12l-6-1 10 14z"
                fill={i % 2 ? "#668e78" : "#769c7e"}
              />
            </g>
          ))}
        </svg>
        <Dragon />
      </div>
    </section>
  );
}
function Puzzle({
  game,
  total,
  parts,
  setParts,
  submit,
  hint,
  feedback,
}: {
  game: Game;
  total: number;
  parts: number[];
  setParts: (v: number[]) => void;
  submit: () => void;
  hint: () => void;
  feedback: string;
}) {
  const sum = parts[0] * 100 + parts[1] * 10 + parts[2];
  return (
    <section className="game">
      <small>
        しずくの谷　／　{game.step + 1} ・ {total}
      </small>
      <h1>小数のかけ橋</h1>
      <p>
        石をつないで、<strong>{decimal(game.target)}</strong>{" "}
        の長さの橋をかけよう。
      </p>
      <div className="river">
        <Dragon />
        <div className="track">
          <div
            className="planks"
            style={{ width: Math.min(100, (100 * sum) / game.target) + "%" }}
          />
        </div>
        <span>⚑</span>
      </div>
      <output>
        いまの長さ　<strong>{decimal(sum)}</strong>
      </output>
      <div className="controls">
        {["1", "0.1", "0.01"].map((v, i) => (
          <fieldset key={v}>
            <legend>{v} の石</legend>
            <div aria-hidden="true" className={"block-model model-" + i} />
            {[-1, 0, 1].map((delta) =>
              delta === 0 ? (
                <strong key={0}> {parts[i]} </strong>
              ) : (
                <button
                  key={delta}
                  aria-label={v + "の石を" + (delta > 0 ? "ふやす" : "へらす")}
                  disabled={parts[i] + delta < 0 || parts[i] + delta > 9}
                  onClick={() =>
                    setParts(parts.map((n, j) => (j === i ? n + delta : n)))
                  }
                >
                  {delta > 0 ? "＋" : "−"}
                </button>
              )
            )}
          </fieldset>
        ))}
      </div>
      <p role="status">{feedback || "　"}</p>
      {game.hint ? (
        <p>１のくらい、小数第１位、小数第２位の順に、石の数を見てみよう。</p>
      ) : (
        <button onClick={hint}>ヒントをみる</button>
      )}{" "}
      <button className="primary" onClick={submit}>
        橋をわたる →
      </button>
    </section>
  );
}
function Report({
  rows,
  create,
}: {
  rows: RecordRow[];
  create: (f: FormData) => void;
}) {
  return (
    <>
      <h1>保護者の部屋</h1>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          create(new FormData(e.currentTarget));
        }}
      >
        <h2>学習者を追加</h2>
        <label>
          ニックネーム
          <input name="name" required maxLength={16} />
        </label>
        <label>
          学校の学年
          <select name="grade">
            {Array.from({ length: 9 }, (_, i) => (
              <option key={i} value={i + 1}>
                {i < 6 ? "小学" + (i + 1) : "中学" + (i - 5)}年
              </option>
            ))}
          </select>
        </label>
        <label>
          アイコン
          <select name="icon">
            {icons.map((x, i) => (
              <option key={i} value={i}>
                {x}
              </option>
            ))}
          </select>
        </label>
        <label>
          本人用の番号
          <input
            name="pin"
            type="password"
            inputMode="numeric"
            pattern="[0-9]{4}"
            placeholder="数字4けた"
            required
          />
        </label>
        <button>登録する</button>
      </form>
      <h2>算数の学習記録</h2>
      {[...new Set(rows.map((r) => r.owner))].map((id) => {
        const history = rows.filter((r) => r.owner === id),
          r = history[0];
        return (
          <article className="record" key={id}>
            <h3>{r.name}　小数のかけ橋</h3>
            <p>
              完了 {history.filter((x) => x.ended).length} 回 ／ 学習時間{" "}
              {history.reduce((n, x) => n + x.seconds, 0)} 秒 ／ 直近の理解度{" "}
              {r.mastery}%　{r.run && r.mastery < 65 ? "・復習候補" : ""}
            </p>
            {history.slice(0, 10).map((x) => (
              <p key={x.run || id}>
                {x.run
                  ? new Date(x.started).toLocaleString("ja-JP") +
                    "　" +
                    (x.ended ? "完了" : "途中") +
                    "　" +
                    x.step +
                    "橋　正解 " +
                    x.correct +
                    "/" +
                    x.answers +
                    "回　" +
                    x.seconds +
                    "秒"
                  : "まだ記録がありません"}
              </p>
            ))}
          </article>
        );
      })}
    </>
  );
}
function App() {
  const [page, P] = React.useState<Page>("profiles"),
    [catalog, C] = React.useState<Catalog>(),
    [list, L] = React.useState<Learner[]>([]),
    [learner, S] = React.useState<Learner>(),
    [rows, R] = React.useState<RecordRow[]>([]),
    [game, G] = React.useState<Game>(),
    [parts, A] = React.useState([0, 0, 0]),
    [elements, E] = React.useState(0),
    [msg, M] = React.useState(""),
    [error, X] = React.useState(""),
    [busy, B] = React.useState(false),
    [still, T] = React.useState(localStorage.motion === "off");
  const lock = React.useRef(false);
  async function work(fn: () => Promise<void>) {
    if (lock.current) return;
    lock.current = true;
    B(true);
    X("");
    try {
      await fn();
    } catch (e) {
      X(e instanceof Error ? e.message : "もういちど試してね");
    } finally {
      B(false);
      lock.current = false;
    }
  }
  async function profiles() {
    L(await api<Learner[]>("profiles"));
    P("profiles");
  }
  async function home() {
    const h = await api<{ elements: string }>("home", { id: learner!.id });
    E(Number(h.elements));
    P("home");
  }
  async function play() {
    G(await api<Game>("play", { id: learner!.id }));
    A([0, 0, 0]);
    M("");
    P("game");
  }
  async function report() {
    R(await api<RecordRow[]>("report"));
    P("report");
  }
  React.useEffect(() => {
    work(async () => {
      C(await (await fetch("/catalog")).json());
      await profiles();
    });
  }, []);
  React.useEffect(() => {
    document.body.classList.toggle("still", still);
    localStorage.motion = still ? "off" : "on";
  }, [still]);
  React.useEffect(() => {
    window.scrollTo(0, 0);
  }, [page]);
  async function exit() {
    await api("logout");
    S(undefined);
    await profiles();
  }
  async function answer() {
    const r = await api<Game>("answer", {
      id: learner!.id,
      step: game!.step,
      run: game!.run,
      parts,
    });
    if (r.correct) {
      G(r);
      P("result");
    } else {
      M(
        parts[0] * 100 + parts[1] * 10 + parts[2] > game!.target
          ? "少し長いみたい。石をへらしてみよう。"
          : "あと少し。石を足してみよう。"
      );
    }
  }
  return (
    <>
      <header>
        <strong>
          <em>◒</em> にじのたまご
        </strong>
        <nav>
          {page === "profiles" ? (
            <button onClick={() => P("parent")}>保護者の部屋</button>
          ) : (
            <button
              disabled={busy}
              onClick={() =>
                work(page === "game" || page === "result" ? home : exit)
              }
            >
              {page === "game" || page === "result"
                ? "島にもどる"
                : "人をえらぶ"}
            </button>
          )}
          <button aria-pressed={still} onClick={() => T(!still)}>
            動きをへらす
          </button>
        </nav>
      </header>
      <fieldset className="screen" disabled={busy}>
        {page === "profiles" && (
          <>
            <Hero title="にじのたまご" caption="NIJI NO TAMAGO">
              <p>きょうは、だれがぼうけんする？</p>
            </Hero>
            <h2>ぼうけんする人</h2>
            <div className="cards">
              {list.map((r) => (
                <button
                  key={r.id}
                  className="card"
                  onClick={() => {
                    S(r);
                    P("select");
                  }}
                >
                  <span className="avatar">{icons[r.icon]}</span>
                  {r.name}
                </button>
              ))}
              {!busy && !list.length && (
                <p>「保護者の部屋」で最初の学習者を登録してください。</p>
              )}
            </div>
          </>
        )}
        {(page === "parent" || page === "select") && (
          <section className="gate">
            <h1>
              {page === "parent"
                ? "保護者の部屋"
                : learner?.name + "のぼうけん"}
            </h1>
            <form
              key={page}
              onSubmit={(e) => {
                e.preventDefault();
                const pin = new FormData(e.currentTarget).get("pin");
                work(async () => {
                  if (page === "parent") {
                    await api("parent", { pin });
                    await report();
                  } else {
                    await api("select", { id: learner!.id, pin });
                    await home();
                  }
                });
              }}
            >
              <label>
                ひみつの番号
                <input
                  name="pin"
                  type="password"
                  inputMode="numeric"
                  autoFocus
                  required
                />
              </label>
              <button className="primary">はいる →</button>
            </form>
          </section>
        )}
        {page === "home" && catalog && (
          <>
            <Hero title="にじの島々" caption="WORLD MAP">
              <p>
                {learner?.name}のぼうけん　✧ {elements} エレメント
              </p>
              <button className="primary" onClick={() => work(play)}>
                小数のかけ橋へ →
              </button>
            </Hero>
            <h2>５つの世界</h2>
            <div className="cards">
              {catalog.worlds.map(([s, w, c], i) => (
                <article className="card" style={{ borderTopColor: c }} key={w}>
                  <Dragon color={c} grown={i === 1 && elements > 0} />
                  <small>{s}</small>
                  <h3>{w}</h3>
                  {i === 1 ? (
                    <button onClick={() => work(play)}>あそぶ →</button>
                  ) : (
                    "じゅんび中"
                  )}
                </article>
              ))}
            </div>
          </>
        )}
        {page === "game" && game && catalog && (
          <Puzzle
            game={game}
            total={catalog.unit.rounds}
            parts={parts}
            setParts={A}
            feedback={msg}
            submit={() => work(answer)}
            hint={() =>
              work(async () =>
                G(
                  await api<Game>("hint", {
                    id: learner!.id,
                    run: game.run,
                    step: game.step,
                  })
                )
              )
            }
          />
        )}
        {page === "result" && game && (
          <section className="gate celebration">
            <Dragon grown={!!game.reward} />
            <small>しずくの谷</small>
            <h1>{game.reward ? "谷をわたりきった！" : "ぴったり！"}</h1>
            <p>
              {game.reward
                ? "青いドラゴンに " + game.reward + " エレメントが届いたよ。"
                : decimal(game.target) + " の橋をわたれたよ。"}
            </p>
            <button onClick={() => work(game.reward ? home : play)}>
              {game.reward ? "島にもどる" : "つぎの橋へ →"}
            </button>
          </section>
        )}
        {page === "report" && (
          <Report
            key={rows.length}
            rows={rows}
            create={(d) =>
              work(async () => {
                await api("create", {
                  name: d.get("name"),
                  grade: Number(d.get("grade")),
                  icon: Number(d.get("icon")),
                  pin: d.get("pin"),
                });
                await report();
              })
            }
          />
        )}
      </fieldset>
      <p id="error" role="alert">
        {error}
      </p>
    </>
  );
}
createRoot(document.getElementById("app")!).render(<App />);
