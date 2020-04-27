require 'spec_helper'

describe package('httpd'), :if => os[:family] == 'ubuntu' do
  it { should_not be_installed }
end

describe user('root') do
  it { should have_uid 0 }
end
