const { exec } = require('child_process');
const fs = require('fs');
const path = require('path');

console.log('🔧 Fixing ChunkLoadError in Frontend...');

// Function to run command
function runCommand(command, description) {
    return new Promise((resolve, reject) => {
        console.log(`📋 ${description}`);
        exec(command, { cwd: __dirname }, (error, stdout, stderr) => {
            if (stdout) console.log(stdout);
            if (stderr) console.log(stderr);
            if (error && !command.includes('pkill')) {
                reject(error);
            } else {
                resolve();
            }
        });
    });
}

// Function to remove directory/file
function removeIfExists(filePath) {
    try {
        if (fs.existsSync(filePath)) {
            fs.rmSync(filePath, { recursive: true, force: true });
            console.log(`✅ Removed: ${filePath}`);
        }
    } catch (error) {
        console.log(`⚠️ Could not remove ${filePath}: ${error.message}`);
    }
}

async function fixChunks() {
    try {
        // Stop any Next.js processes
        await runCommand('pkill -f "next dev"', 'Stopping Next.js processes');
        
        // Clean cache and dependencies
        console.log('📦 Cleaning cache and dependencies...');
        removeIfExists('.next');
        removeIfExists('node_modules');
        removeIfExists('.npm');
        removeIfExists('package-lock.json');
        
        // Install dependencies
        await runCommand('npm install', 'Installing dependencies');
        
        // Build project
        await runCommand('npm run build', 'Building project');
        
        // Clean build for development
        console.log('🧹 Cleaning build for development...');
        removeIfExists('.next');
        
        console.log('✅ Frontend fixed! Ready to start development server');
        console.log('🚀 Run: npm run dev');
        console.log('🌐 Access at: http://localhost:3000');
        
    } catch (error) {
        console.error('❌ Error:', error.message);
        process.exit(1);
    }
}

fixChunks();