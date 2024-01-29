// Fetch courses data from the backend API
const CourseIPSchool = "http://10.2.130.98:9090/courses";
const CourseIPSHome = "http://192.168.1.10:9090/courses";
//const ReserveIpSchool = `http://10.2.130.98:9090/course/${courseId}/reserve/${seat}`;
//const ReserveIpHome = `http://192.168.1.10:9090/course/${courseId}/reserve/${seat}`;
const token = sessionStorage.getItem('token');

fetch('http://10.2.130.45:9090/admincourses', {
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

  const deletebtn = document.createElement('button');
  deletebtn.className = "Delete course"
  deletebtn.textContent = "Delete course"
  deletebtn.onclick = () => showDeleteCoursePopUp(course.ID);
  courseCard.appendChild(deletebtn);
  return courseCard;
}
const showDeleteCoursePopUp = (id) => {
    const div = document.getElementById('reserve-modal')

    div.innerHTML = `
        <div id="modal">
            <div id="content">
                <p>Delete course?</p>
                <div class="yes-no">
                   <button onclick="deleteCourse(${id})">YES</button>
                   <button onclick="closeModal()">NO</button>
                </div>
            </div>
        </div>
    `
}
function deleteCourse(courseId) {
    const token = sessionStorage.getItem('token');
    if (!token) {
      console.log('User is not logged in.'); // Replace with desired error handling
      return;
    }
    
      fetch(`http://10.2.130.45:9090/course/${courseId}/delete`, {
      method: 'DELETE',
      headers: {
        'token': `${token}`
      }
    })
      .then(response => {
        if (response.status === 200) {
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

  function CreateCourse(date , description, title){
    const token = sessionStorage.getItem("token")
    if (!token) {
        console.log('User is not logged in.'); // Replace with desired error handling
        return;
      }
    fetch(`http://10.2.130.45:9090/addcourse`, {
      method: 'PUT',
      headers: {
        'token': `${token}`
      },
      body: JSON.stringify({
        date: new Date(date).toISOString(),
        description: description,
        title: title,
      })
    })
      .then(response => {
        if (response.status === 200) {
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
  const showPopup = () => {
    window.scrollTo(0,0)
    document.getElementById('reserve-modal').innerHTML = `
    <div id="modal">
  <div id="content">
    <span id="close-modal" onclick="closeModal()">&times;</span>
    <h1>Create Course</h1>
    <div>
      <label for="course-title">Title:</label>
      <input type="text" id="course-title" placeholder="Title">
    </div>
    <div>
      <label for="course-date">Date:</label>
      <input type="datetime-local" id="course-date" placeholder="Date and Time">
    </div>
    <div>
      <label for="course-description">Description:</label>
      <textarea id="course-description" placeholder="Description"></textarea>
    </div>
    <div id="seat-container"></div>
    <button id="submit-btn" onclick="submitForm()">Submit</button>
  </div>
</div>

    `     
   
}
function submitForm(){
const title = document.getElementById('course-title').value;
  const date = document.getElementById('course-date').value;
  const description = document.getElementById('course-description').value;

  // Call the method to handle the submitted values
  
  CreateCourse(date,description,title);
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
