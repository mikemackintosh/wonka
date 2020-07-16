require 'serverspec'

set :backend, :exec

describe package('httpd'), :if => os[:family] == 'ubuntu' do
  it { should_not be_installed }
end

describe user('root') do
  it { should have_uid 0 }
end

describe user('splug') do
  it { should_not exist }
end
