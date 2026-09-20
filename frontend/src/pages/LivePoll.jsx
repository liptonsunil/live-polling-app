import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';

export default function LivePoll() {
  const { id } = useParams();
  const [poll, setPoll] = useState(null);

  // 1. Fetch initial poll data
  useEffect(() => {
    const fetchPoll = async () => {
      try {
        const res = await axios.get(`http://localhost:8080/api/polls/${id}`);
        setPoll(res.data);
      } catch (error) {
        console.error("Error fetching poll:", error);
      }
    };
    fetchPoll();
  }, [id]);

  // 2. Establish the WebSocket connection
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws/polls/${id}`);
    ws.onmessage = (event) => {
      const updatedPoll = JSON.parse(event.data);
      setPoll(updatedPoll);
    };
    return () => ws.close();
  }, [id]);

  // 3. Handle casting a vote with the new Voter ID
  const handleVote = async (optionId) => {
    let voterId = localStorage.getItem("voterId");
    if (!voterId) {
      voterId = "voter_" + Math.random().toString(36).substr(2, 9);
      localStorage.setItem("voterId", voterId);
    }

    try {
      await axios.post(`http://localhost:8080/api/polls/${id}/vote`, {
        option_id: optionId,
        voter_id: voterId
      });
    } catch (error) {
      if (error.response && error.response.status === 403) {
        alert("You have already voted on this poll!");
      } else {
        console.error("Error casting vote:", error);
      }
    }
  };

  if (!poll) return <p>Loading poll...</p>;

  return (
    <div style={{ padding: '20px', maxWidth: '500px', margin: '0 auto' }}>
      <h2>{poll.question}</h2>
      <p style={{ color: 'gray' }}>Share this link: {window.location.href}</p>
      <div style={{ marginTop: '20px' }}>
        {poll.options.map((opt) => (
          <div key={opt.id} style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '10px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>{opt.text} — <strong>{opt.votes} votes</strong></span>
            <button onClick={() => handleVote(opt.id)} style={{ padding: '5px 10px', cursor: 'pointer' }}>Vote</button>
          </div>
        ))}
      </div>
    </div>
  );
}
