require 'spec_helper'

describe package('httpd'), :if => os[:family] == 'ubuntu' do
  it { should_not be_installed }
end

describe user('splug') do
  it { should exist }
  it { should have_uid 500 }
  it { should have_home_directory '/home/splug' }
  it { should have_login_shell '/bin/bash' }
end
