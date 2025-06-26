const fs = require('fs');
const path = require('path');

console.log('🧹 Starting cleanup process...');

// Function to remove directory/file safely
function removeIfExists(filePath) {
    try {
        const fullPath = path.resolve(__dirname, filePath);
        if (fs.existsSync(fullPath)) {
            fs.rmSync(fullPath, { recursive: true, force: true });
            console.log(`✅ Removed: ${filePath}`);
            return true;
        } else {
            console.log(`ℹ️ Not found: ${filePath}`);
            return false;
        }
    } catch (error) {
        console.log(`⚠️ Could not remove ${filePath}: ${error.message}`);
        return false;
    }
}

// Clean up files
console.log('📦 Cleaning cache and dependencies...');

const removed = [];
if (removeIfExists('.next')) removed.push('.next');
if (removeIfExists('node_modules')) removed.push('node_modules');
if (removeIfExists('.npm')) removed.push('.npm');
if (removeIfExists('package-lock.json')) removed.push('package-lock.json');

console.log(`\n🗑️ Cleanup completed. Removed: ${removed.length > 0 ? removed.join(', ') : 'nothing (already clean)'}`);
console.log('\n📋 Next steps:');
console.log('1. Run: npm install');
console.log('2. Run: npm run build');  
console.log('3. Run: npm run dev');
console.log('4. Access: http://localhost:3000');
console.log('\n✅ Ready for manual execution of npm commands!');