package main

import (
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
	case "addschd":
		for i := 0; i < system.Blastoff; i++ {
			sch.Add <- nil
		}
	case "add":
		sch.Add <- nil
	case "blastoff":
		sch.Blastoff <- system.Blastoff
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

const Tmpl = `<html>

<head>
   <style>
      html {
         background-color: black;
         font-family: Arial, Helvetica, sans-serif;
      }

      canvas {
         background-image: radial-gradient(circle at top right, rgba(0, 97, 114, 0.1), transparent),
            radial-gradient(circle at bottom left, rgba(0, 97, 114, 0.2), transparent),
            radial-gradient(circle at top, rgba(0, 162, 255, 0.11), transparent),
            radial-gradient(circle at bottom, rgba(0, 162, 255, 0.11), transparent);
         z-index: -8888;
         position: fixed;
         opacity: 0.5;
         width: 104%;
         height: 104%;
         margin-top: -50px;
         margin-left: -20px;
      }

      .center {
         margin: auto;
         width: 50%;
         padding: 10px;
      }

      .center-more {
         margin: auto;
         text-align: center;
         padding: 10px;
         color: white;
      }

      #viewport {
         margin-top: 10px;
         margin-bottom: 10px;
      }

      p {
         color: white;
         text-align: center;
         width: 50%;
         margin: auto;
         padding: 20px;
      }
   </style>
</head>
<body>
   <div>
      <canvas></canvas>
      <div class="center-more">
         <h1>Space Elevator</h1>
         <button onclick="addOne()">Add [1] Passenger [Using NC Algo] Via [Scheduler]</button>
         <button onclick="addschd()">Add [200] Passengers [Using NC Algo] Via [Scheduler]</button>
         <button onclick="blastOff()">Add [200] Passengers To [Each Ship]!</button>
      </div>
      <div style="font-size: 16px; text-align: center;">
         <p>Idle: 👩‍🚀&nbsp;&nbsp;&nbsp;&nbsp;Up: 🚀&nbsp;&nbsp;&nbsp;&nbsp;Down: 🛸</p>
         <pre id="viewport"></pre>
      </div>
      <div style="font-size: 16px; text-align: center;">
        <p>
            Nearest Car (NC): Elevator calls are assigned to the elevator best placed to answer that call according to three criteria that are used to compute a figure of suitability (FS) for each elevator. (1) If an elevator is moving towards a call, and the call is in the same direction, FS = (N + 2) - d, where N is one less than the number of floors in the building, and d is the distance in floors between the elevator and the passenger call. (2) If the elevator is moving towards the call, but the call is in the opposite direction, FS = (N + 1) - d.  (3) If the elevator is moving away from the point of call, FS = 1. The elevator with the highest FS for each call is sent to answer it. The search for the "nearest car" is performed continuously until each call is serviced.
        </p>
      </div>
      <div class="center" style="font-size: 14px; color: white; ">
         <div style="display: inline-block";>
            <pre id="stats"></pre>
         </div>
      </div>
   </div>

   <script>

      paint();
      stats();

      setInterval(function () {
         paint();
      }, 1000);
      
      setInterval(function () {
         stats();
      }, 1000);


      function stats() {
         send("http://localhost:8989/cmd/stats")
         .then(function(response) {
            let vp = document.getElementById( 'stats' );
            vp.innerHTML = response; 
         })
         .catch(error => console.error(error));
      }

      function paint() {
         send("http://localhost:8989/cmd/paint")
         .then(function(response) {
            let vp = document.getElementById( 'viewport' );
            vp.innerHTML = response; 
         })
         .catch(error => console.error(error));
      }

      function addschd() {
         send("http://localhost:8989/cmd/addschd")
            .catch(error => console.error(error));
      }

      function addOne() {
         send("http://localhost:8989/cmd/add")
            .catch(error => console.error(error));
      }

      function blastOff() {
         send("http://localhost:8989/cmd/blastoff")
            .catch(error => console.error(error));
      }


      function send(url = "", data = {}) {
         return fetch(url, {
            method: "GET", // *GET, POST, PUT, DELETE, etc.
            mode: "cors", // no-cors, cors, *same-origin
            cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
            credentials: "same-origin", // include, *same-origin, omit
            redirect: "follow", // manual, *follow, error
            referrer: "no-referrer", // no-referrer, *client
         })
            .then(response => response.text());
      }
   </script>


   <script>
            const STAR_COUNT = (window.innerWidth + window.innerHeight) / 8,
               STAR_SIZE = 3,
               STAR_MIN_SCALE = 0.2,
               OVERFLOW_THRESHOLD = 50;

            const canvas = document.querySelector('canvas'),
               context = canvas.getContext('2d');

            let scale = 1, // device pixel ratio
               width,
               height;

            let stars = [];

            let pointerX,
               pointerY;

            let velocity = { x: 0, y: 0, tx: 0, ty: 0, z: 0.0005 };

            let touchInput = false;

            generate();
            resize();
            step();

            let mouseMoveSpeed = 150;
            window.onresize = resize;
            window.onmousemove = onMouseMove;
            window.ontouchmove = onTouchMove;
            window.ontouchend = onMouseLeave;
            document.onmouseleave = onMouseLeave;

            function generate() {

               for (let i = 0; i < STAR_COUNT; i++) {
                  stars.push({
                     x: 0,
                     y: 0,
                     z: STAR_MIN_SCALE + Math.random() * (1 - STAR_MIN_SCALE)
                  });
               }

            }

            function placeStar(star) {

               star.x = Math.random() * width;
               star.y = Math.random() * height;

            }

            function recycleStar(star) {

               let direction = 'z';

               let vx = Math.abs(velocity.tx), vy = Math.abs(velocity.ty);

               if (vx > 1 && vy > 1) {
                  let axis;

                  if (vx > vy) {
                     axis = Math.random() < Math.abs(velocity.x) / (vx + vy) ? 'h' : 'v';
                  }
                  else {
                     axis = Math.random() < Math.abs(velocity.y) / (vx + vy) ? 'v' : 'h';
                  }

                  if (axis === 'h') {
                     direction = velocity.x > 0 ? 'l' : 'r';
                  }
                  else {
                     direction = velocity.y > 0 ? 't' : 'b';
                  }
               }

               star.z = STAR_MIN_SCALE + Math.random() * (1 - STAR_MIN_SCALE);

               if (direction === 'z') {
                  star.z = 0.1;
                  star.x = Math.random() * width;
                  star.y = Math.random() * height;
               }
               else if (direction === 'l') {
                  star.x = -STAR_SIZE;
                  star.y = height * Math.random();
               }
               else if (direction === 'r') {
                  star.x = width + STAR_SIZE;
                  star.y = height * Math.random();
               }
               else if (direction === 't') {
                  star.x = width * Math.random();
                  star.y = -STAR_SIZE;
               }
               else if (direction === 'b') {
                  star.x = width * Math.random();
                  star.y = height + STAR_SIZE;
               }

            }

            function resize() {

               scale = window.devicePixelRatio || 1;

               width = window.innerWidth * scale;
               height = window.innerHeight * scale;

               canvas.width = width;
               canvas.height = height;

               stars.forEach(placeStar);

            }

            function step() {

               context.clearRect(0, 0, width, height);

               update();
               render();

               requestAnimationFrame(step);

            }

            function update() {

               velocity.tx *= 0.95;
               velocity.ty *= 0.95;

               velocity.x += (velocity.tx - velocity.x) * 0.7;
               velocity.y += (velocity.ty - velocity.y) * 0.7;

               stars.forEach(function (star) {

                  star.x += velocity.x * star.z;
                  star.y += velocity.y * star.z;

                  star.x += (star.x - width / 2) * velocity.z * star.z;
                  star.y += (star.y - height / 2) * velocity.z * star.z;
                  star.z += velocity.z;

                  // recycle when out of bounds
                  if (star.x < -OVERFLOW_THRESHOLD || star.x > width + OVERFLOW_THRESHOLD || star.y < -OVERFLOW_THRESHOLD || star.y > height + OVERFLOW_THRESHOLD) {
                     recycleStar(star);
                  }

               });

            }

            function render() {

               stars.forEach(function (star) {

                  context.beginPath();
                  context.lineCap = 'round';
                  context.lineWidth = STAR_SIZE * star.z * scale;
                  context.strokeStyle = 'rgba(255,255,255,' + (0.5 + 0.5 * Math.random()) + ')';

                  context.beginPath();
                  context.moveTo(star.x, star.y);

                  var tailX = velocity.x * 2,
                     tailY = velocity.y * 2;

                  // stroke() wont work on an invisible line
                  if (Math.abs(tailX) < 0.1) tailX = 0.5;
                  if (Math.abs(tailY) < 0.1) tailY = 0.5;

                  context.lineTo(star.x + tailX, star.y + tailY);

                  context.stroke();

               });

            }

            function movePointer(x, y) {

               if (typeof pointerX === 'number' && typeof pointerY === 'number') {

                  let ox = x - pointerX,
                     oy = y - pointerY;

                  velocity.tx = velocity.x + (ox / mouseMoveSpeed * scale) * (touchInput ? -1 : 1);
                  velocity.ty = velocity.y + (oy / mouseMoveSpeed * scale) * (touchInput ? -1 : 1);

               }

               pointerX = x;
               pointerY = y;

            }

            function onMouseMove(event) {

               touchInput = false;

               movePointer(event.clientX, event.clientY);

            }

            function onTouchMove(event) {

               touchInput = true;

               movePointer(event.touches[0].clientX, event.touches[0].clientY, true);

               event.preventDefault();

            }

            function onMouseLeave() {

               pointerX = null;
               pointerY = null;

            }
   </script>
</body>

</html>`
