import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

export default function CreatePoll() {
  const [question, setQuestion] = useState('');
  const [options, setOptions] = useState([{ id: '1', text: '', votes: 0 }, { id: '2', text: '', votes: 0 }]);
  const navigate = useNavigate();

  // Kick the user to the login screen if they don't have a token
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) navigate('/login');
  }, [navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    const creator = localStorage.getItem('username'); // Use actual logged-in user

    try {
      const response = await axios.post('https://polling-backend-c9oo.onrender.com/api/polls', {
        creator: creator,
        question,
        options
      }, {
        // This header unlocks the Go AuthMiddleware
        headers: { Authorization: `Bearer ${token}` } 
      });
      navigate(`/poll/${response.data.id}`);
    } catch (error) {
      if (error.response?.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
      } else {
        console.error("Error creating poll:", error);
      }
    }
  };

  const updateOption = (index, text) => {
    const newOptions = [...options];
    newOptions[index].text = text;
    setOptions(newOptions);
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/login');
  };

  return (
    <div style={{ padding: '20px', maxWidth: '500px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Create a Live Poll</h2>
        <button onClick={handleLogout} style={{ padding: '5px 10px', cursor: 'pointer' }}>Logout</button>
      </div>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '15px' }}>
          <label style={{ display: 'block', marginBottom: '5px' }}>Question:</label>
          <input type="text" value={question} onChange={(e) => setQuestion(e.target.value)} required style={{ width: '100%', padding: '8px' }} />
        </div>
        
        {options.map((opt, index) => (
          <div key={index} style={{ marginBottom: '10px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Option {index + 1}:</label>
            <input type="text" value={opt.text} onChange={(e) => updateOption(index, e.target.value)} required style={{ width: '100%', padding: '8px' }} />
          </div>
        ))}
        
        <button type="submit" style={{ padding: '10px 20px', cursor: 'pointer', marginTop: '10px' }}>
          Create & Share Link
        </button>
      </form>
    </div>
  );
}
