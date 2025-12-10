import React, { useState } from 'react';
import { Card } from '../components/ui/Card';
import { LLM_MODELS, PROTOCOLS } from '../constants';
import { Settings, Save, Database, FileText, Check, X } from 'lucide-react';

const ModelConfig: React.FC = () => {
  const [selectedModel, setSelectedModel] = useState(LLM_MODELS[0].id);
  const [temperature, setTemperature] = useState(0.2);
  const [maxTokens, setMaxTokens] = useState(2048);
  const [ragEnabled, setRagEnabled] = useState(true);

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-slate-900">Model & RAG Configuration</h2>
        <p className="text-slate-500 mt-1">Manage LLM parameters and clinical knowledge base sources.</p>
      </div>

      {/* Model Selection */}
      <Card title="LLM Configuration" action={
        <button className="flex items-center px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700 text-sm font-medium">
          <Save className="h-4 w-4 mr-2" /> Save Changes
        </button>
      }>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">Primary Scoring Model</label>
              <select 
                value={selectedModel}
                onChange={(e) => setSelectedModel(e.target.value)}
                className="block w-full rounded-md border-slate-300 shadow-sm focus:border-teal-500 focus:ring-teal-500 sm:text-sm p-2.5 border"
              >
                {LLM_MODELS.map(model => (
                  <option key={model.id} value={model.id}>{model.name}</option>
                ))}
              </select>
              <div className="mt-2 flex flex-wrap gap-2">
                {LLM_MODELS.find(m => m.id === selectedModel)?.tags.map(tag => (
                  <span key={tag} className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-800">
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">Embedding Model</label>
              <select className="block w-full rounded-md border-slate-300 shadow-sm focus:border-teal-500 focus:ring-teal-500 sm:text-sm p-2.5 border">
                <option>text-embedding-3-large (OpenAI)</option>
                <option>BAAI/bge-m3 (Local)</option>
              </select>
            </div>
          </div>

          <div className="space-y-6 bg-slate-50 p-6 rounded-lg border border-slate-100">
            <div>
              <div className="flex justify-between">
                <label className="block text-sm font-medium text-slate-700">Temperature</label>
                <span className="text-sm text-slate-500">{temperature}</span>
              </div>
              <input 
                type="range" 
                min="0" 
                max="1" 
                step="0.1"
                value={temperature}
                onChange={(e) => setTemperature(parseFloat(e.target.value))}
                className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer mt-2"
              />
            </div>

            <div>
              <div className="flex justify-between">
                <label className="block text-sm font-medium text-slate-700">Max Output Tokens</label>
                <span className="text-sm text-slate-500">{maxTokens}</span>
              </div>
              <input 
                type="range" 
                min="512" 
                max="8192" 
                step="512"
                value={maxTokens}
                onChange={(e) => setMaxTokens(parseInt(e.target.value))}
                className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer mt-2"
              />
            </div>

             <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">System Instruction Preview</label>
                <textarea 
                  rows={4}
                  className="block w-full rounded-md border-slate-300 shadow-sm focus:border-teal-500 focus:ring-teal-500 sm:text-xs text-slate-600 p-2 border bg-white"
                  readOnly
                  value="You are an expert medical examiner. Assess the student's answer based on the provided clinical protocols. Identify clinical accuracy, completeness, and adherence to guidelines. Provide a score from 0-100 and structured feedback."
                />
             </div>
          </div>
        </div>
      </Card>

      {/* RAG Knowledge Base */}
      <Card 
        title="Knowledge Base Sources" 
        subtitle="Manage clinical protocols used for Retrieval Augmented Generation"
        action={
           <div className="flex items-center space-x-3">
              <span className="text-sm text-slate-600">RAG for Scoring</span>
              <button 
                onClick={() => setRagEnabled(!ragEnabled)}
                className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none ${ragEnabled ? 'bg-teal-600' : 'bg-slate-200'}`}
              >
                <span className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${ragEnabled ? 'translate-x-5' : 'translate-x-0'}`} />
              </button>
           </div>
        }
      >
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-slate-200">
            <thead className="bg-slate-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Protocol Name</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Type</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Last Updated</th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Status</th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-slate-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-slate-200">
              {PROTOCOLS.map((protocol) => (
                <tr key={protocol.id}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <FileText className="h-5 w-5 text-slate-400 mr-2" />
                      <div className="text-sm font-medium text-slate-900">{protocol.name}</div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500">
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-800">
                      {protocol.type}
                    </span>
                    <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-800">
                      {protocol.language}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500">{protocol.last_updated}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {protocol.status === 'indexed' ? (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Indexed
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                        Indexing...
                      </span>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-teal-600 hover:text-teal-900 mr-4">Preview Chunks</button>
                    {protocol.active ? (
                        <button className="text-slate-400 hover:text-red-600"><Check className="h-4 w-4" /></button>
                    ) : (
                        <button className="text-slate-400 hover:text-green-600"><X className="h-4 w-4" /></button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
};

export default ModelConfig;
