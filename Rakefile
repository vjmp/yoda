desc 'Run unit tests.'
task :test do
    sh "go test ./..."
end

task :hawk do
    sh "hawk .go /// go test ./..."
end

task :default => :test
