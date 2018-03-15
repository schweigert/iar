require 'benchmark'
require 'colorize'
require 'pry'

class AntBase
  attr_accessor :position_length, :position_height

  def initialize(map)
    @map = map
    @random = Random.new

    init
  end

  def think
    @map.move(ant: self) do
      @position_height = (@position_height + @random.rand(3) - 1) % @map.height
      @position_length = (@position_length + @random.rand(3) - 1) % @map.length
    end
  end

  private

  def init
    # initial position
    @position_height = @random.rand(@map.height)
    @position_length = @random.rand(@map.length)


    @map.add(ant: self)
  end
end

class DeadAnt < AntBase; end

class Ant < AntBase; end

class Map
  attr_accessor :length, :height

  def initialize(length: 50, height: 50, ants: 10, dead_ants: 100, ants_per_thread: 20)
    @length = length
    @height = height

    @specimens = ants
    @dead_specimens = dead_ants

    @ants_per_thread = ants_per_thread

    init
  end

  def move(ant: nil)
    remove(ant: ant)
    yield
    add(ant: ant)
  end

  def remove(ant: nil)
    unless ant.nil?
      @grid[ant.position_length][ant.position_height].delete ant
    end
  end

  def add(ant: nil)
    unless ant.nil?
      @grid[ant.position_length][ant.position_height] << ant
    end
  end

  def puts
    @grid.each do |line|
      line.each do |cell|
        case cell.last
        when Ant
          print '@ '.blue
        when DeadAnt
          print '# '.red
        else
          print '  '
        end
      end

      print "\n"
    end
  end

  def interate(interations)
    @groups = @ants.each_slice(@ants_per_thread).to_a
    @threads = []

    @groups.each do |group|
      @threads << Thread.new(group, interations) do |ants, n|
        n.times do
          ants.each do |ant|
            ant.think
          end
        end
      end
    end

    @threads.each(&:join)
  end

  private

  def init
    # init grid
    @grid ||= []

    for i in 0...@length
      @grid[i] ||= []

      for j in 0...@height
        @grid[i][j] ||= []
      end
    end

    # create ants
    @ants = []
    @specimens.times { @ants << Ant.new(self) }

    # create dead_ants
    @dead_ants = []
    @dead_specimens.times { @dead_ants << DeadAnt.new(self) }
  end
end

# Start a new map
map = Map.new(ants: 100)

loop do
  system 'clear'
  puts Benchmark.measure { map.interate(1) }
  map.puts
end
