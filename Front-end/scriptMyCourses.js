
const logoutIP = "http://10.2.130.45:9090/logout" ;


function logout(){
    const token = sessionStorage.getItem('token');
    if (!token) {
      console.log('User is not logged in.'); // Replace with desired error handling
      return;
    }
    fetch(`http://10.2.130.45:9090/logout`, {
        method: 'POST',
        headers: {
          'token': `${token}`
        }
      })
        .then(response => {
          if (response.status === 200) {
            console.log('logout successfuly');
            window.location.assign("/loginPage.html")
          } else if (response.status === 401) {
            console.log('Unauthorized: you are not logged'); // Replace with desired error handling
          }
        })
        .catch(error => console.error('Error:', error));

}
function courses(){
    window.location.assign("/course.html")
}
const token = sessionStorage.getItem('token');
// Fetch courses data from the backend API
fetch('http://10.2.130.45:9090/mycourses', {
  method: 'GET',
  headers: {
    'token': token
  }
})
  .then(response => response.json())
  .then(data => {
    const coursesContainer = document.getElementById('courses-container');
    data.forEach(course => {
      const courseCard = createCourseCard(course);
      coursesContainer.appendChild(courseCard);
    });
  })
  .catch(error => console.error('Error:', error));

// Create a course card
function createCourseCard(course) {
  const courseCard = document.createElement('div');
  courseCard.className = 'course-card';

  const title = document.createElement('div');
  title.className = 'course-title';
  title.textContent = course.Title;
  courseCard.appendChild(title);

  const description = document.createElement('div');
  description.className = 'course-description';
  description.textContent = course.Description;
  courseCard.appendChild(description);

  const details = document.createElement('div');
  details.className = 'course-details';
  details.innerHTML = `
    <div>Date: ${new Date(course.Date).toLocaleString()}</div>
    <div>Lector: ${course.Lector}</div>
    <div>Participants: ${course.Participants ? course.Participants.length : 0}/15</div>
  `;
  courseCard.appendChild(details);

  const reserveButton = document.createElement('button');

  let isParticipating = false
  if(course.Participants){
    for(let participant of course.Participants){
        if(participant.Username === sessionStorage.getItem("username")){
            isParticipating = true;
            break;
        }
      }
  }

  reserveButton.className = isParticipating ? 'seat-reserved' : 'seat-available';
  reserveButton.textContent = isParticipating ? 'Remove reservation' : 'Reserve Seat';
  console.log(course)
  reserveButton.onclick = isParticipating ? () => showUnreservePopup(course.ID) : () => showPopup(course);
  courseCard.appendChild(reserveButton);

  return courseCard;
}

const showUnreservePopup = (id) => {
    const div = document.getElementById('reserve-modal')

    div.innerHTML = `
        <div id="modal">
            <div id="content">
                <p>Cancel the reservation?</p>
                <div class="yes-no">
                   <button onclick="unreserveSeat(${id})">YES</button>
                   <button onclick="closeModal()">NO</button>
                </div>
            </div>
        </div>
    `
}

// Reserve a seat for a course
function reserveSeat(courseId, seat) {
  const token = sessionStorage.getItem('token');
  if (!token) {
    console.log('User is not logged in.'); // Replace with desired error handling
    return;
  }

  fetch(`http://10.2.130.45:9090/course/${courseId}/reserve/${seat}`, {
    method: 'POST',
    headers: {
      'token': `${token}`
    }
  })
    .then(response => {
      if (response.status === 201) {
        console.log('Seat reserved successfully!');
        window.location.reload()
      } else if (response.status === 401) {
        console.log('Unauthorized: User is not logged in.'); // Replace with desired error handling
      } else {
        console.log('Seat reservation failed.'); // Replace with desired error handling
      }
    })
    .catch(error => console.error('Error:', error));
}

function unreserveSeat(courseId) {
    const token = sessionStorage.getItem('token');
    if (!token) {
      console.log('User is not logged in.'); // Replace with desired error handling
      return;
    }
  
    fetch(`http://10.2.130.45:9090/course/${courseId}/unreserved`, {
      method: 'PUT',
      headers: {
        'token': `${token}`
      }
    })
      .then(response => {
        if (response.status === 201) {
          console.log('Unreserved from course successfully!');
          window.location.reload()
        } else if (response.status === 401) {
          console.log('Unauthorized: User is not logged in.'); // Replace with desired error handling
        } else {
            console.log(response.status)
          console.log('Seat unreservation failed.'); // Replace with desired error handling
        }
      })
      .catch(error => console.error('Error:', error));
  }

const showPopup = (data) => {
    window.scrollTo(0,0)
    let takenSeats = []
    if(data.Participants){
        for(let participant of data.Participants){
            takenSeats.push(participant.Seat)
        }
    }

    document.getElementById('reserve-modal').innerHTML = `
    <div id="modal">
      <div id="content">
        <span id="close-modal" onclick="closeModal()">&times;</span>
        <h1>Title: ${data.Title}</h1>
        <p>Lector: ${data.Lector}</p>
        <div id="seat-container"></div>
        </div>
      </div>
    `     
    for(let i = 1; i <= 15; i++){
        if(takenSeats.includes(i)){
            document.getElementById('seat-container').innerHTML += `
            <button class="seat-btn" disabled>Seat Taken</button>
        `
        } else {
            document.getElementById('seat-container').innerHTML += `
            <button class="seat-btn-${i}" onclick="confirm(${data.ID},${i})">Seat ${i}</button>
        `
        }
        //
    }
}

const confirm = (id, seat) => {
    const btn = document.getElementsByClassName(`seat-btn-${seat}`)[0]
    const div = document.getElementById('content');

    btn.classList.add('active-btn')
    div.innerHTML += `
        <div id="confirmation">
            <p>Reserve seat ${seat}?</p>
            <div class="yes-no">
                <button onclick="reserveSeat(${id}, ${seat})">YES</button>
                <button onclick="closeModal()">NO</button>
            </div>
        </div>
    `
}

const closeModal = () => {
    document.getElementById('reserve-modal').innerHTML = ''
}
const logoutBtn = document.getElementById('logout-btn');

function logout(){
    const token = sessionStorage.getItem('token');
    if (!token) {
      console.log('User is not logged in.'); // Replace with desired error handling
      return;
    }
    fetch(`http://10.2.130.45:9090/logout`, {
        method: 'POST',
        headers: {
          'token': `${token}`
        }
      })
        .then(response => {
          if (response.status === 200) {
            console.log('logout successfuly');
            window.location.assign("/loginPage.html")
          } else if (response.status === 401) {
            console.log('Unauthorized: you are not logged'); // Replace with desired error handling
          }
        })
        .catch(error => console.error('Error:', error));

}
function myCourses(){
    window.location.assign("/myCourses.html")
}