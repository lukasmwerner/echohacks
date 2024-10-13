Project Title: EchoHacks
Overview
EchoHacks is a cutting-edge SSH-based application that empowers hackathon participants by providing a collaborative platform for sharing and evaluating ideas. Inspired by Reddit, EchoHacks allows users to submit hackathon concepts and express their support through upvotes and downvotes. The application features a terminal interface built with the Bubble Tea framework, ensuring a smooth and engaging user experience with real-time feedback on idea popularity.

Features
Terminal Interface: A sleek and intuitive command-line interface designed for easy interaction and navigation.
Idea Submission: Users can submit their hackathon ideas, including titles, ranks, and usernames.
Voting System: Community members can upvote or downvote ideas, reflecting their support or disapproval.
Color-Coded Feedback: Ideas are displayed with colors indicating their popularity—green for positive (upvotes) and red for negative (downvotes).
Real-Time Updates: All changes in votes and rankings are reflected instantly, keeping users informed and engaged.
SQLite Database: Utilizes SQLite for efficient storage and retrieval of ideas and their associated votes.
Technology Stack
Programming Language: Go
Framework: Bubble Tea for terminal-based UI
Database: modernc.org/sqlite for data management
SSH Handling: Charmbracelet's Wish framework for managing SSH sessions
Styling: Lipgloss for enhanced terminal aesthetics
Installation
Clone the repository:
bash
Copy code
git clone https://github.com/yourusername/echohacks.git
Navigate to the project directory:
bash
Copy code
cd echohacks
Install dependencies using Go modules:
bash
Copy code
go mod tidy
Ensure SQLite is installed on your machine and create the database:
bash
Copy code
touch app.db
Run the application:
bash
Copy code
go run main.go
Usage
Start the application and connect via SSH to access the terminal interface.
Users can create new posts to share their hackathon ideas.
Participants can upvote or downvote ideas, with visual feedback reflecting their overall popularity.
Contribution
We welcome contributions to enhance the functionality and user experience of EchoHacks. Please feel free to open issues or submit pull requests.

License
This project is licensed under the MIT License. See the LICENSE file for details.

Acknowledgments
Special thanks to the Charmbracelet team for their powerful libraries that made building EchoHacks a rewarding experience.
