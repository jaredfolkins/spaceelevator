package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaredfolkins/spaceelevator/system"
	"github.com/pkg/browser"
)

var router *mux.Router
var sch *system.Scheduler

type appState struct {
	Time      string           `json:"time"`
	Graph     [][]string       `json:"graph"`
	Statuses  []*system.Status `json:"statuses"`
	Floors    int              `json:"floors"`
	Elevators int              `json:"elevators"`
}

func main() {
	sch = system.NewScheduler(system.Floors, system.Elevators)
	sch.Run()

	router = mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/cmd/{cmd}", Cmd).Methods("GET")

	go func() {
		time.Sleep(1 * time.Second)
		browser.OpenURL("http://localhost:8989")
	}()

	log.Fatal(http.ListenAndServe(":8989", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(Tmpl))
}

func Cmd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cmd := vars["cmd"]
	switch cmd {
	case "stats":
		ch := make(chan []*system.Status)
		sch.Status <- ch
		st := <-ch
		close(ch)
		res := stats(st)
		t := time.Now()
		stime := fmt.Sprintln(t.Format("2006-01-02 15:04:05")) + "\n"
		tmpl := "ElevatorID\t\tDestinationFloor\t\tCurrentFloor\t\tDirection\t\tPickupQueue\t\tDropoffQueue\n"
		w.Write([]byte(stime))
		w.Write([]byte(tmpl))
		w.Write([]byte(res))
	case "paint":
		ch := make(chan *system.Paint)
		sch.Paint <- ch
		p := <-ch
		close(ch)
		res := paint(p)
		w.Write([]byte(res))
	case "state":
		statusCh := make(chan []*system.Status)
		sch.Status <- statusCh
		st := <-statusCh
		close(statusCh)

		paintCh := make(chan *system.Paint)
		sch.Paint <- paintCh
		p := <-paintCh
		close(paintCh)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(appState{
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			Graph:     p.Graph,
			Statuses:  st,
			Floors:    system.Floors,
			Elevators: system.Elevators,
		})
	case "addschd":
		for i := 0; i < system.Blastoff; i++ {
			sch.Add <- nil
		}
	case "add":
		sch.Add <- nil
	case "blastoff":
		sch.Blastoff <- system.Blastoff
	default:
		http.NotFound(w, r)
	}
	return
}

func paint(p *system.Paint) string {
	var res string
	for i := len(p.Graph) - 1; i >= 0; i-- {
		for _, l := range p.Graph[i] {
			res = res + l
		}
		res = res + "\n"
	}
	return res
}

func stats(st []*system.Status) string {
	var res string
	tmpl := "%d\t\t\t%d\t\t\t\t%d\t\t\t%s\t\t\t%d\t\t\t%d\n"
	for _, s := range st {
		res = res + fmt.Sprintf(tmpl, s.ElevatorID, s.DestinationFloor, s.CurrentFloor, s.Direction, s.PickupTotal, s.DropoffTotal)
	}
	res = res + "\n"
	return res
}

const Tmpl = `<!doctype html>
<html lang="en">

<head>
   <meta charset="utf-8">
   <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
   <title>Space Elevator</title>
   <link rel="icon" href="data:,">
   <style>
      :root {
         color-scheme: dark;
         --void: #030406;
         --panel: rgba(9, 10, 13, 0.82);
         --panel-strong: rgba(14, 15, 18, 0.94);
         --paper: #f5f1e6;
         --muted: #b8b1a1;
         --line: rgba(245, 241, 230, 0.2);
         --line-strong: rgba(245, 241, 230, 0.34);
         --red: #e3342f;
         --amber: #e7b75f;
         --blue: #6aa7ff;
         --cyan: #62d6ce;
         --green: #80d084;
      }

      * {
         box-sizing: border-box;
      }

      html,
      body {
         min-height: 100%;
         margin: 0;
      }

      body {
         overflow-x: hidden;
         background:
            linear-gradient(180deg, rgba(255, 255, 255, 0.08), transparent 18rem),
            linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px),
            linear-gradient(180deg, rgba(255, 255, 255, 0.025) 1px, transparent 1px),
            var(--void);
         background-size: auto, 64px 64px, 64px 64px, auto;
         color: var(--paper);
         font-family: "Avenir Next", "Inter", "Segoe UI", Arial, Helvetica, sans-serif;
         letter-spacing: 0;
      }

      body::before {
         position: fixed;
         inset: 0;
         z-index: -2;
         pointer-events: none;
         content: "";
         background:
            linear-gradient(90deg, transparent 0 11%, rgba(255, 255, 255, 0.06) 11% 11.16%, transparent 11.16% 88.84%, rgba(255, 255, 255, 0.06) 88.84% 89%, transparent 89%),
            linear-gradient(180deg, rgba(0, 0, 0, 0.34), rgba(0, 0, 0, 0.72));
      }

      #starfield {
         position: fixed;
         inset: 0;
         z-index: -3;
         width: 100%;
         height: 100%;
         opacity: 0.72;
         pointer-events: none;
      }

      .main-shell {
         width: min(1180px, calc(100% - 32px));
         margin: 0 auto;
         padding: clamp(16px, 3vw, 36px) 0 28px;
         display: grid;
         gap: 16px;
      }

      .command-deck,
      .simulation-panel,
      .fleet-panel,
      .algorithm-panel {
         border: 1px solid var(--line);
         border-radius: 8px;
         background: var(--panel);
         box-shadow: 0 20px 80px rgba(0, 0, 0, 0.36);
         backdrop-filter: blur(14px);
      }

      .command-deck {
         min-height: 124px;
         display: grid;
         grid-template-columns: minmax(0, 1fr) auto;
         align-items: stretch;
         gap: 14px;
         padding: clamp(14px, 2.4vw, 24px);
      }

      .identity {
         display: flex;
         align-items: center;
         min-width: 0;
         gap: clamp(12px, 2vw, 20px);
      }

      .monolith-mark {
         flex: 0 0 clamp(34px, 5vw, 56px);
         align-self: stretch;
         min-height: 92px;
         border: 1px solid rgba(255, 255, 255, 0.16);
         border-radius: 2px;
         background:
            linear-gradient(105deg, rgba(255, 255, 255, 0.16), transparent 28%),
            linear-gradient(180deg, #050507, #17181c 58%, #020203);
         box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.85), 0 16px 40px rgba(0, 0, 0, 0.44);
      }

      .eyebrow {
         margin: 0 0 6px;
         color: var(--amber);
         font-size: clamp(0.68rem, 1.3vw, 0.78rem);
         font-weight: 800;
         text-transform: uppercase;
      }

      h1,
      h2 {
         margin: 0;
         padding: 0.04em 0;
         font-weight: 850;
         line-height: 1.25;
      }

      h1 {
         max-width: 10ch;
         font-size: clamp(2.1rem, 8.8vw, 6.8rem);
         text-transform: uppercase;
      }

      h2 {
         font-size: clamp(1.08rem, 2.4vw, 1.55rem);
         text-transform: uppercase;
      }

      .hal-status {
         min-width: min(100%, 238px);
         padding: 14px;
         border: 1px solid var(--line-strong);
         border-radius: 6px;
         background:
            linear-gradient(180deg, rgba(255, 255, 255, 0.08), transparent),
            rgba(0, 0, 0, 0.52);
         display: grid;
         grid-template-columns: 58px minmax(0, 1fr);
         align-items: center;
         gap: 12px;
      }

      .hal-lens {
         width: 56px;
         height: 56px;
         border-radius: 999px;
         background:
            radial-gradient(circle at 50% 50%, #fff6ef 0 5%, #ffb36f 6% 13%, var(--red) 14% 36%, #5e090b 37% 62%, #090203 63%),
            #270304;
         box-shadow: 0 0 24px rgba(227, 52, 47, 0.55), inset 0 0 16px rgba(255, 255, 255, 0.16);
      }

      .hal-copy {
         min-width: 0;
         color: var(--muted);
         font-size: 0.76rem;
         font-weight: 800;
         text-transform: uppercase;
      }

      .hal-copy strong {
         display: block;
         color: var(--paper);
         font-size: 1rem;
      }

      .system-mode {
         margin-top: 5px;
         color: var(--green);
      }

      .actions {
         display: grid;
         grid-template-columns: repeat(3, minmax(0, 1fr));
         gap: 10px;
      }

      button {
         min-height: 46px;
         border: 1px solid var(--line-strong);
         border-radius: 6px;
         background: rgba(245, 241, 230, 0.08);
         color: var(--paper);
         font: inherit;
         font-weight: 800;
         cursor: pointer;
         touch-action: manipulation;
      }

      button:hover,
      button:focus-visible {
         border-color: rgba(245, 241, 230, 0.7);
         background: rgba(245, 241, 230, 0.15);
         outline: none;
      }

      button:disabled {
         cursor: wait;
         opacity: 0.56;
      }

      .command-button {
         display: grid;
         grid-template-columns: auto minmax(0, 1fr);
         align-items: center;
         justify-items: start;
         gap: 10px;
         padding: 0 14px;
         text-align: left;
      }

      .button-mark {
         width: 28px;
         height: 28px;
         border: 1px solid var(--line-strong);
         border-radius: 999px;
         display: grid;
         place-items: center;
         color: var(--amber);
         font-size: 1.1rem;
         line-height: 1;
      }

      .panel-head {
         display: flex;
         align-items: center;
         justify-content: space-between;
         gap: 12px;
         padding: 16px;
         border-bottom: 1px solid var(--line);
      }

      .segmented {
         display: inline-grid;
         grid-template-columns: repeat(2, minmax(86px, 1fr));
         padding: 3px;
         border: 1px solid var(--line);
         border-radius: 7px;
         background: rgba(0, 0, 0, 0.38);
      }

      .segmented button {
         min-height: 36px;
         border: 0;
         border-radius: 4px;
         background: transparent;
         color: var(--muted);
         font-size: 0.82rem;
      }

      .segmented button[aria-pressed="true"] {
         background: var(--paper);
         color: #07080a;
      }

      .telemetry-strip {
         display: grid;
         grid-template-columns: repeat(4, minmax(0, 1fr));
         gap: 1px;
         border-bottom: 1px solid var(--line);
         background: var(--line);
      }

      .telemetry-pill {
         min-width: 0;
         padding: 13px 16px;
         background: rgba(0, 0, 0, 0.42);
      }

      .telemetry-pill span {
         display: block;
         color: var(--muted);
         font-size: 0.7rem;
         font-weight: 800;
         text-transform: uppercase;
      }

      .telemetry-pill strong {
         display: block;
         margin-top: 3px;
         overflow-wrap: anywhere;
         font-size: clamp(1.12rem, 2.6vw, 1.6rem);
      }

      .shaft-frame {
         --cols: 100;
         --rows: 16;
         --cell-size: 6px;
         --cell-gap: 1px;
         --target-cell-gap: clamp(0.35px, 0.14vw, 1.6px);
         --floor-font-size: 6px;
         display: grid;
         grid-template-columns: clamp(18px, 4.5vw, 42px) clamp(18px, 5.4vw, 72px) minmax(0, 1fr);
         gap: clamp(4px, 1vw, 10px);
         padding: clamp(10px, 1.7vw, 16px);
         overflow: hidden;
      }

      .shaft-frame.is-detail {
         --target-cell-gap: clamp(0.6px, 0.2vw, 2.2px);
      }

      .shaft-frame.is-detail .activity-row {
         border-color: rgba(245, 241, 230, 0.26);
      }

      .floor-scale,
      .floor-activity {
         display: grid;
         align-self: start;
         grid-template-rows: repeat(var(--rows, 16), var(--cell-size));
         gap: var(--cell-gap);
      }

      .floor-scale {
         color: var(--muted);
         font-size: var(--floor-font-size);
         font-weight: 800;
         text-align: right;
         text-transform: uppercase;
      }

      .floor-scale span {
         height: var(--cell-size);
         min-height: 0;
         display: flex;
         align-items: center;
         justify-content: flex-end;
         overflow: hidden;
         line-height: 1;
      }

      .floor-activity {
         min-width: 0;
      }

      .activity-row {
         height: var(--cell-size);
         min-height: 0;
         position: relative;
         overflow: hidden;
         border: 1px solid rgba(245, 241, 230, 0.11);
         border-radius: 2px;
         background: rgba(245, 241, 230, 0.055);
      }

      .activity-row span {
         position: absolute;
         left: 0;
         width: 0;
         min-width: 0;
         transition: width 160ms ease;
      }

      .activity-pickup {
         top: 0;
         height: 50%;
         background: var(--cyan);
         box-shadow: 0 0 8px rgba(98, 214, 206, 0.48);
      }

      .activity-dropoff {
         bottom: 0;
         height: 50%;
         background: var(--amber);
         box-shadow: 0 0 8px rgba(231, 183, 95, 0.46);
      }

      .activity-row.is-empty {
         opacity: 0.48;
      }

      .grid-scroll {
         min-width: 0;
         overflow: hidden;
      }

      .lift-grid {
         display: grid;
         grid-template-columns: repeat(var(--cols), var(--cell-size));
         grid-auto-rows: var(--cell-size);
         gap: var(--cell-gap);
         width: max-content;
      }

      .lift-cell {
         width: var(--cell-size);
         height: var(--cell-size);
         min-width: 0;
         min-height: 0;
         border-radius: min(2px, calc(var(--cell-size) / 4));
         background: rgba(245, 241, 230, 0.07);
      }

      .shaft-frame.is-detail .lift-cell {
         outline: 1px solid rgba(245, 241, 230, 0.16);
         outline-offset: -1px;
      }

      .shaft-frame.is-detail .lift-cell:nth-child(10n) {
         outline-color: rgba(231, 183, 95, 0.58);
      }

      .shaft-frame.is-detail .lift-cell:nth-child(100n + 1) {
         outline-color: rgba(98, 214, 206, 0.62);
      }

      .lift-cell.is-idle {
         background: var(--paper);
         box-shadow: 0 0 9px rgba(245, 241, 230, 0.52);
      }

      .shaft-frame.is-detail .lift-cell.is-idle {
         box-shadow: 0 0 12px rgba(245, 241, 230, 0.72);
      }

      .lift-cell.is-up {
         background: linear-gradient(180deg, #ffffff, var(--blue) 38%, #194f8a);
         box-shadow: 0 0 10px rgba(106, 167, 255, 0.7);
      }

      .shaft-frame.is-detail .lift-cell.is-up {
         box-shadow: 0 0 14px rgba(106, 167, 255, 0.9);
      }

      .lift-cell.is-down {
         background: linear-gradient(180deg, #fff0c7, var(--amber) 42%, #9a4d15);
         box-shadow: 0 0 10px rgba(231, 183, 95, 0.68);
      }

      .shaft-frame.is-detail .lift-cell.is-down {
         box-shadow: 0 0 14px rgba(231, 183, 95, 0.88);
      }

      .legend {
         display: flex;
         flex-wrap: wrap;
         gap: 12px;
         padding: 0 16px 16px;
         color: var(--muted);
         font-size: 0.78rem;
         font-weight: 800;
         text-transform: uppercase;
      }

      .legend-item {
         display: inline-flex;
         align-items: center;
         gap: 7px;
      }

      .legend-swatch {
         width: 15px;
         height: 15px;
         border-radius: 2px;
         background: rgba(245, 241, 230, 0.07);
      }

      .legend-swatch.is-idle {
         background: var(--paper);
      }

      .legend-swatch.is-up {
         background: var(--blue);
      }

      .legend-swatch.is-down {
         background: var(--amber);
      }

      .legend-swatch.is-pickup {
         background: var(--cyan);
      }

      .legend-swatch.is-dropoff {
         background: var(--amber);
      }

      .fleet-panel {
         overflow: hidden;
      }

      .table-scroll {
         max-height: 340px;
         overflow: auto;
      }

      .stats-table {
         width: 100%;
         min-width: 680px;
         border-collapse: collapse;
         font-size: 0.84rem;
      }

      .stats-table th,
      .stats-table td {
         padding: 10px 12px;
         border-bottom: 1px solid rgba(245, 241, 230, 0.1);
         text-align: left;
         white-space: nowrap;
      }

      .stats-table th {
         position: sticky;
         top: 0;
         z-index: 1;
         background: var(--panel-strong);
         color: var(--amber);
         font-size: 0.72rem;
         text-transform: uppercase;
      }

      .stats-table tr[data-direction="up"] td:first-child {
         color: var(--blue);
      }

      .stats-table tr[data-direction="down"] td:first-child {
         color: var(--amber);
      }

      .stats-table tr[data-direction="idle"] td:first-child {
         color: var(--paper);
      }

      .algorithm-panel {
         padding: 16px;
      }

      .algorithm-panel p {
         margin: 0;
         max-width: 76ch;
         color: var(--muted);
         font-size: clamp(0.96rem, 1.9vw, 1.08rem);
         line-height: 1.72;
      }

      .sr-only {
         position: absolute;
         width: 1px;
         height: 1px;
         padding: 0;
         margin: -1px;
         overflow: hidden;
         clip: rect(0, 0, 0, 0);
         white-space: nowrap;
         border: 0;
      }

      @media (max-width: 760px) {
         .main-shell {
            width: min(100% - 20px, 1180px);
            padding-top: 10px;
         }

         .command-deck {
            grid-template-columns: 1fr;
            min-height: 0;
         }

         .identity {
            align-items: stretch;
         }

         .monolith-mark {
            min-height: 82px;
         }

         h1 {
            max-width: 9ch;
            font-size: clamp(2.15rem, 14vw, 3.9rem);
         }

         .hal-status {
            grid-template-columns: 48px minmax(0, 1fr);
            min-width: 0;
         }

         .hal-lens {
            width: 46px;
            height: 46px;
         }

         .actions {
            grid-template-columns: 1fr;
         }

         .panel-head {
            display: grid;
            align-items: stretch;
         }

         .segmented {
            width: 100%;
         }

         .telemetry-strip {
            grid-template-columns: repeat(2, minmax(0, 1fr));
         }

         .shaft-frame {
            grid-template-columns: 24px 34px minmax(0, 1fr);
            gap: 7px;
            padding: 12px;
         }

         .legend {
            padding: 0 12px 12px;
            gap: 9px;
         }

         .table-scroll {
            max-height: 310px;
         }

         .stats-table {
            min-width: 560px;
            font-size: 0.78rem;
         }

         .stats-table th,
         .stats-table td {
            padding: 9px 10px;
         }
      }
   </style>
</head>

<body>
   <canvas id="starfield" aria-hidden="true"></canvas>

   <main class="main-shell">
      <section class="command-deck" aria-labelledby="app-title">
         <div class="identity">
            <div class="monolith-mark" aria-hidden="true"></div>
            <div>
               <p class="eyebrow">Jupiter Mission // Elevator Control</p>
               <h1 id="app-title">Space Elevator</h1>
            </div>
         </div>
         <div class="hal-status" aria-label="Mission computer status">
            <span class="hal-lens" aria-hidden="true"></span>
            <div class="hal-copy">
               <strong>HAL 9000</strong>
               <span id="renderedAt">Awaiting signal</span>
               <div id="systemMode" class="system-mode" aria-live="polite">ONLINE</div>
            </div>
         </div>
      </section>

      <section class="actions" aria-label="Passenger controls">
         <button class="command-button" type="button" data-cmd="add">
            <span class="button-mark" aria-hidden="true">+</span>
            <span>Add 1 Passenger</span>
         </button>
         <button class="command-button" type="button" data-cmd="addschd">
            <span class="button-mark" aria-hidden="true">200</span>
            <span>Schedule 200</span>
         </button>
         <button class="command-button" type="button" data-cmd="blastoff">
            <span class="button-mark" aria-hidden="true">ALL</span>
            <span>Fill Every Car</span>
         </button>
      </section>

      <section class="simulation-panel" aria-labelledby="matrix-title">
         <div class="panel-head">
            <div>
               <p class="eyebrow">Orbital Lift Matrix</p>
               <h2 id="matrix-title">Transit Grid</h2>
            </div>
            <div class="segmented" aria-label="Grid density">
               <button type="button" data-density="overview" aria-pressed="true">Overview</button>
               <button type="button" data-density="detail" aria-pressed="false">Detail</button>
            </div>
         </div>
         <div class="telemetry-strip" aria-label="Fleet telemetry">
            <div class="telemetry-pill">
               <span>Active Cars</span>
               <strong id="activeCars">0</strong>
            </div>
            <div class="telemetry-pill">
               <span>Pickup Queue</span>
               <strong id="pickupQueue">0</strong>
            </div>
            <div class="telemetry-pill">
               <span>Dropoff Queue</span>
               <strong id="dropoffQueue">0</strong>
            </div>
            <div class="telemetry-pill">
               <span>Fleet Size</span>
               <strong id="fleetSize">0</strong>
            </div>
         </div>
         <div class="shaft-frame">
            <div id="floorScale" class="floor-scale" aria-hidden="true"></div>
            <div id="floorActivity" class="floor-activity" aria-label="Pickup and dropoff activity by floor"></div>
            <div class="grid-scroll">
               <div id="viewport" class="lift-grid" role="grid" aria-label="Elevator positions"></div>
            </div>
         </div>
         <div class="legend" aria-label="Transit legend">
            <span class="legend-item"><span class="legend-swatch is-idle" aria-hidden="true"></span>Idle</span>
            <span class="legend-item"><span class="legend-swatch is-up" aria-hidden="true"></span>Ascending</span>
            <span class="legend-item"><span class="legend-swatch is-down" aria-hidden="true"></span>Descending</span>
            <span class="legend-item"><span class="legend-swatch" aria-hidden="true"></span>Open Shaft</span>
            <span class="legend-item"><span class="legend-swatch is-pickup" aria-hidden="true"></span>Pickup Calls</span>
            <span class="legend-item"><span class="legend-swatch is-dropoff" aria-hidden="true"></span>Dropoff Targets</span>
         </div>
      </section>

      <section class="fleet-panel" aria-labelledby="fleet-title">
         <div class="panel-head">
            <div>
               <p class="eyebrow">Discovery Fleet Telemetry</p>
               <h2 id="fleet-title">Car Status</h2>
            </div>
         </div>
         <div class="table-scroll">
            <table class="stats-table">
               <thead>
                  <tr>
                     <th scope="col">Car</th>
                     <th scope="col">Direction</th>
                     <th scope="col">Current</th>
                     <th scope="col">Destination</th>
                     <th scope="col">Pickup</th>
                     <th scope="col">Dropoff</th>
                  </tr>
               </thead>
               <tbody id="stats"></tbody>
            </table>
         </div>
      </section>

      <section class="algorithm-panel" aria-labelledby="algorithm-title">
         <p class="eyebrow">Elevator Algorithm</p>
         <h2 id="algorithm-title">Nearest Car Algorithm</h2>
         <p>
            Elevator calls are assigned to the car best placed to answer that call according to three figure-of-suitability checks. A car moving toward a same-direction call scores highest, a car moving toward an opposite-direction call scores next, and a car moving away from the call scores lowest. The highest-scoring car is assigned while the search runs continuously until each call is serviced.
         </p>
      </section>
   </main>

   <script>
      const els = {
         activeCars: document.getElementById('activeCars'),
         pickupQueue: document.getElementById('pickupQueue'),
         dropoffQueue: document.getElementById('dropoffQueue'),
         fleetSize: document.getElementById('fleetSize'),
         floorActivity: document.getElementById('floorActivity'),
         floorScale: document.getElementById('floorScale'),
         renderedAt: document.getElementById('renderedAt'),
         shaftFrame: document.querySelector('.shaft-frame'),
         stats: document.getElementById('stats'),
         systemMode: document.getElementById('systemMode'),
         viewport: document.getElementById('viewport'),
         densityButtons: document.querySelectorAll('[data-density]'),
         commandButtons: document.querySelectorAll('[data-cmd]')
      };

      let renderedFloorCount = 0;
      let density = window.localStorage.getItem('spaceelevator-density') || 'overview';
      let shaftLayoutFrame = 0;

      function setDensity(nextDensity) {
         density = nextDensity;
         window.localStorage.setItem('spaceelevator-density', density);
         els.viewport.classList.toggle('is-detail', density === 'detail');
         els.shaftFrame.classList.toggle('is-detail', density === 'detail');
         scheduleShaftLayout();
         els.densityButtons.forEach(function(button) {
            button.setAttribute('aria-pressed', button.dataset.density === density ? 'true' : 'false');
         });
      }

      function scheduleShaftLayout() {
         window.cancelAnimationFrame(shaftLayoutFrame);
         shaftLayoutFrame = window.requestAnimationFrame(syncShaftLayout);
      }

      function px(value, fallback) {
         const parsed = Number.parseFloat(value);
         return Number.isFinite(parsed) ? parsed : fallback;
      }

      function syncShaftLayout() {
         const style = window.getComputedStyle(els.shaftFrame);
         const cols = Math.max(1, Number.parseInt(style.getPropertyValue('--cols'), 10) || 100);
         const padding = px(style.paddingLeft, 0) + px(style.paddingRight, 0);
         const frameGap = px(style.columnGap, px(style.gap, 0));
         const labelWidth = els.floorScale.getBoundingClientRect().width;
         const activityWidth = els.floorActivity.getBoundingClientRect().width;
         const availableWidth = els.shaftFrame.clientWidth - padding - labelWidth - activityWidth - (frameGap * 2);
         const targetGap = px(style.getPropertyValue('--target-cell-gap'), 1);
         let cellGap = Math.max(0.2, targetGap);
         let cellSize = (availableWidth - ((cols - 1) * cellGap)) / cols;

         if (cellSize < 1.25) {
            cellGap = Math.max(0.2, Math.min(cellGap, availableWidth / (cols * 6)));
            cellSize = (availableWidth - ((cols - 1) * cellGap)) / cols;
         }

         cellSize = Math.max(0.8, cellSize);
         const floorFontSize = Math.max(1, Math.min(9, cellSize * 0.88));

         els.shaftFrame.style.setProperty('--cell-size', cellSize.toFixed(3) + 'px');
         els.shaftFrame.style.setProperty('--cell-gap', cellGap.toFixed(3) + 'px');
         els.shaftFrame.style.setProperty('--floor-font-size', floorFontSize.toFixed(3) + 'px');
      }

      function classifyCell(symbol) {
         if (symbol === '🌎') {
            return { name: 'empty', label: 'open shaft' };
         }
         if (symbol === '🚀') {
            return { name: 'up', label: 'ascending' };
         }
         if (symbol === '🛸') {
            return { name: 'down', label: 'descending' };
         }
         return { name: 'idle', label: 'idle' };
      }

      function renderFloorScale(rowCount) {
         if (renderedFloorCount === rowCount) {
            return;
         }
         renderedFloorCount = rowCount;
         els.floorScale.style.setProperty('--rows', String(rowCount));
         const fragment = document.createDocumentFragment();
         for (let i = rowCount; i >= 1; i--) {
            const label = document.createElement('span');
            label.textContent = 'F' + i;
            fragment.appendChild(label);
         }
         els.floorScale.replaceChildren(fragment);
      }

      function aggregateFloorActivity(statuses, rowCount) {
         const pickup = Array(rowCount).fill(0);
         const dropoff = Array(rowCount).fill(0);

         statuses.forEach(function(status) {
            if (Array.isArray(status.PickupFloors)) {
               status.PickupFloors.forEach(function(count, floor) {
                  if (floor < rowCount) {
                     pickup[floor] += count;
                  }
               });
            }
            if (Array.isArray(status.DropoffFloors)) {
               status.DropoffFloors.forEach(function(count, floor) {
                  if (floor < rowCount) {
                     dropoff[floor] += count;
                  }
               });
            }
         });

         return { pickup: pickup, dropoff: dropoff };
      }

      function renderFloorActivity(statuses, rowCount) {
         els.floorActivity.style.setProperty('--rows', String(rowCount));

         const activity = aggregateFloorActivity(statuses, rowCount);
         const maxCount = Math.max(1, ...activity.pickup, ...activity.dropoff);
         const fragment = document.createDocumentFragment();

         for (let floor = rowCount - 1; floor >= 0; floor--) {
            const pickupCount = activity.pickup[floor];
            const dropoffCount = activity.dropoff[floor];
            const row = document.createElement('span');
            const pickupFill = document.createElement('span');
            const dropoffFill = document.createElement('span');

            row.className = 'activity-row' + (pickupCount + dropoffCount === 0 ? ' is-empty' : '');
            row.setAttribute('role', 'meter');
            row.setAttribute('aria-label', 'F' + (floor + 1) + ': ' + pickupCount + ' pickup calls, ' + dropoffCount + ' dropoff targets');
            row.setAttribute('aria-valuemin', '0');
            row.setAttribute('aria-valuemax', String(maxCount));
            row.setAttribute('aria-valuenow', String(Math.max(pickupCount, dropoffCount)));

            pickupFill.className = 'activity-pickup';
            pickupFill.style.width = (pickupCount / maxCount * 100) + '%';
            dropoffFill.className = 'activity-dropoff';
            dropoffFill.style.width = (dropoffCount / maxCount * 100) + '%';

            row.appendChild(pickupFill);
            row.appendChild(dropoffFill);
            fragment.appendChild(row);
         }

         els.floorActivity.replaceChildren(fragment);
      }

      function renderGrid(graph) {
         const rows = graph.slice().reverse();
         const rowCount = rows.length;
         const colCount = rowCount > 0 ? rows[0].length : 0;

         renderFloorScale(rowCount);
         els.shaftFrame.style.setProperty('--rows', String(rowCount));
         els.shaftFrame.style.setProperty('--cols', String(colCount));
         els.viewport.style.setProperty('--rows', String(rowCount));
         els.viewport.style.setProperty('--cols', String(colCount));
         els.viewport.setAttribute('aria-rowcount', String(rowCount));
         els.viewport.setAttribute('aria-colcount', String(colCount));
         scheduleShaftLayout();

         const fragment = document.createDocumentFragment();
         rows.forEach(function(row, rowIndex) {
            const floor = rowCount - rowIndex;
            row.forEach(function(symbol, colIndex) {
               const state = classifyCell(symbol);
               const cell = document.createElement('span');
               cell.className = 'lift-cell is-' + state.name;
               cell.setAttribute('role', 'gridcell');
               cell.setAttribute('aria-rowindex', String(rowIndex + 1));
               cell.setAttribute('aria-colindex', String(colIndex + 1));
               cell.setAttribute('aria-label', 'Car ' + (colIndex + 1) + ', floor ' + floor + ', ' + state.label);
               fragment.appendChild(cell);
            });
         });
         els.viewport.replaceChildren(fragment);
      }

      function appendCell(row, value) {
         const cell = document.createElement('td');
         cell.textContent = value;
         row.appendChild(cell);
      }

      function renderStats(statuses) {
         const fragment = document.createDocumentFragment();
         let activeCars = 0;
         let pickupQueue = 0;
         let dropoffQueue = 0;

         statuses.forEach(function(status) {
            pickupQueue += status.PickupTotal;
            dropoffQueue += status.DropoffTotal;
            if (status.Direction !== 'idle' || status.PickupTotal > 0 || status.DropoffTotal > 0) {
               activeCars++;
            }

            const row = document.createElement('tr');
            row.dataset.direction = status.Direction;
            appendCell(row, 'Car ' + (status.ElevatorID + 1));
            appendCell(row, status.Direction.toUpperCase());
            appendCell(row, 'F' + (status.CurrentFloor + 1));
            appendCell(row, 'F' + (status.DestinationFloor + 1));
            appendCell(row, String(status.PickupTotal));
            appendCell(row, String(status.DropoffTotal));
            fragment.appendChild(row);
         });

         els.activeCars.textContent = String(activeCars);
         els.pickupQueue.textContent = String(pickupQueue);
         els.dropoffQueue.textContent = String(dropoffQueue);
         els.fleetSize.textContent = String(statuses.length);
         els.stats.replaceChildren(fragment);
      }

      function renderState(state) {
         els.renderedAt.textContent = state.time;
         els.systemMode.textContent = 'ONLINE';
         els.systemMode.style.color = 'var(--green)';
         renderGrid(state.graph || []);
         renderFloorActivity(state.statuses || [], (state.graph || []).length);
         renderStats(state.statuses || []);
      }

      function showError(error) {
         console.error(error);
         els.systemMode.textContent = 'SIGNAL LOST';
         els.systemMode.style.color = 'var(--red)';
      }

      function refresh() {
         return fetch('/cmd/state', { cache: 'no-cache' })
            .then(function(response) {
               if (!response.ok) {
                  throw new Error('State request failed with ' + response.status);
               }
               return response.json();
            })
            .then(renderState)
            .catch(showError);
      }

      function runCommand(command, button) {
         button.disabled = true;
         fetch('/cmd/' + command, { cache: 'no-cache' })
            .then(function(response) {
               if (!response.ok) {
                  throw new Error(command + ' failed with ' + response.status);
               }
               return refresh();
            })
            .catch(showError)
            .finally(function() {
               button.disabled = false;
            });
      }

      els.commandButtons.forEach(function(button) {
         button.addEventListener('click', function() {
            runCommand(button.dataset.cmd, button);
         });
      });

      els.densityButtons.forEach(function(button) {
         button.addEventListener('click', function() {
            setDensity(button.dataset.density);
         });
      });

      if ('ResizeObserver' in window) {
         const resizeObserver = new ResizeObserver(scheduleShaftLayout);
         resizeObserver.observe(els.shaftFrame);
      } else {
         window.addEventListener('resize', scheduleShaftLayout);
      }

      setDensity(density);
      refresh();
      window.setInterval(refresh, 1000);
   </script>

   <script>
      (function() {
         const canvas = document.getElementById('starfield');
         const context = canvas.getContext('2d');
         const starCount = Math.min(320, Math.max(120, Math.floor((window.innerWidth + window.innerHeight) / 6)));
         const starSize = 2.4;
         const minScale = 0.18;
         const overflow = 50;

         let scale = 1;
         let width = 0;
         let height = 0;
         let pointerX = null;
         let pointerY = null;
         let velocity = { x: 0, y: 0, tx: 0, ty: 0, z: 0.00034 };
         let stars = [];

         function generate() {
            stars = [];
            for (let i = 0; i < starCount; i++) {
               stars.push({
                  x: 0,
                  y: 0,
                  z: minScale + Math.random() * (1 - minScale)
               });
            }
         }

         function placeStar(star) {
            star.x = Math.random() * width;
            star.y = Math.random() * height;
         }

         function recycleStar(star) {
            star.z = minScale + Math.random() * (1 - minScale);
            star.x = Math.random() * width;
            star.y = Math.random() * height;
         }

         function resize() {
            scale = window.devicePixelRatio || 1;
            width = Math.floor(window.innerWidth * scale);
            height = Math.floor(window.innerHeight * scale);
            canvas.width = width;
            canvas.height = height;
            stars.forEach(placeStar);
         }

         function update() {
            velocity.tx *= 0.965;
            velocity.ty *= 0.965;
            velocity.x += (velocity.tx - velocity.x) * 0.08;
            velocity.y += (velocity.ty - velocity.y) * 0.08;

            stars.forEach(function(star) {
               star.x += velocity.x * star.z;
               star.y += velocity.y * star.z;
               star.x += (star.x - width / 2) * velocity.z * star.z;
               star.y += (star.y - height / 2) * velocity.z * star.z;
               star.z += velocity.z;

               if (star.x < -overflow || star.x > width + overflow || star.y < -overflow || star.y > height + overflow) {
                  recycleStar(star);
               }
            });
         }

         function render() {
            context.clearRect(0, 0, width, height);
            stars.forEach(function(star) {
               const alpha = 0.22 + star.z * 0.58;
               const tailX = Math.abs(velocity.x) < 0.1 ? 0.35 : velocity.x * 1.7;
               const tailY = Math.abs(velocity.y) < 0.1 ? 0.35 : velocity.y * 1.7;

               context.beginPath();
               context.lineCap = 'round';
               context.lineWidth = starSize * star.z * scale;
               context.strokeStyle = 'rgba(245, 241, 230, ' + alpha + ')';
               context.moveTo(star.x, star.y);
               context.lineTo(star.x + tailX, star.y + tailY);
               context.stroke();
            });
         }

         function step() {
            update();
            render();
            window.requestAnimationFrame(step);
         }

         function movePointer(x, y) {
            if (typeof pointerX === 'number' && typeof pointerY === 'number') {
               velocity.tx = velocity.x + (x - pointerX) / 170 * scale;
               velocity.ty = velocity.y + (y - pointerY) / 170 * scale;
            }
            pointerX = x;
            pointerY = y;
         }

         function resetPointer() {
            pointerX = null;
            pointerY = null;
         }

         window.addEventListener('resize', resize);
         window.addEventListener('pointermove', function(event) {
            if (event.pointerType === 'touch') {
               return;
            }
            movePointer(event.clientX, event.clientY);
         }, { passive: true });
         document.addEventListener('mouseleave', resetPointer);

         generate();
         resize();
         step();
      })();
   </script>
</body>

</html>`
