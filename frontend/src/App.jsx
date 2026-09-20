import { BrowserRouter, Routes, Route } from 'react-router-dom';
import CreatePoll from './pages/CreatePoll';
import LivePoll from './pages/LivePoll';
import Login from './pages/Login';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<CreatePoll />} />
        <Route path="/poll/:id" element={<LivePoll />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
