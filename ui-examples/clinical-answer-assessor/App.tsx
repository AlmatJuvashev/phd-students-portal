
import React, { useState } from 'react';
import { Routes, Route, HashRouter } from 'react-router-dom';
import Layout from './components/ui/Layout';
import Dashboard from './pages/Dashboard';
import ModelConfig from './pages/ModelConfig';
import ScoreAudit from './pages/ScoreAudit';
import TrainingGround from './pages/TrainingGround';
import TrainingSession from './pages/TrainingSession';
import GradingQueue from './pages/GradingQueue';
import Landing from './pages/Landing';
import { UserRole } from './types';

const App: React.FC = () => {
  // Global State for Demo
  const [role, setRole] = useState<UserRole>('admin');
  const [currentModel] = useState<string>('Qwen 2 72B');

  return (
    <HashRouter>
      <Routes>
        <Route path="/" element={<Landing />} />
        
        {/* App Routes wrapped in Layout */}
        <Route element={<Layout role={role} setRole={setRole} currentModel={currentModel} />}>
          <Route path="/dashboard" element={<Dashboard role={role} />} />
          <Route path="/config" element={<ModelConfig />} />
          <Route path="/audit" element={<ScoreAudit />} />
          <Route path="/training" element={<TrainingGround />} />
          <Route path="/training/:topicId" element={<TrainingSession />} />
          <Route path="/grading" element={<GradingQueue />} />
        </Route>
      </Routes>
    </HashRouter>
  );
};

export default App;
