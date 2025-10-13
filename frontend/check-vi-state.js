// Paste this into browser console to check VI module state

const journeyState = JSON.parse(sessionStorage.getItem('journeyState') || '{}');
const viNodes = [
  'VI1_post_defense_overview',
  'VI2_library_deposits', 
  'VI3_ncste_state_registration',
  'VI_attestation_file'
];

console.log('=== Module VI State ===');
viNodes.forEach(nodeId => {
  console.log(`${nodeId}: ${journeyState[nodeId] || 'not set'}`);
});

const allDone = viNodes.every(id => journeyState[id] === 'done');
console.log('\nAll VI nodes done?', allDone);

// Clear old module guard flags
console.log('\n=== Clearing module_guard flags ===');
for (let i = 0; i < sessionStorage.length; i++) {
  const key = sessionStorage.key(i);
  if (key && key.startsWith('module_guard_')) {
    console.log('Removing:', key);
    sessionStorage.removeItem(key);
  }
}
console.log('Done! Refresh page to see clean state.');
