require 'pry'
require 'thread'
require 'colorize'
require 'benchmark'

class AntBase
  attr_accessor :position_length, :position_height

  def initialize(map)
    @map = map
    @random = Random.new

    init
  end

  def think
  end

  private

  def init
    # initial position
    @position_height = @random.rand(@map.height)
    @position_length = @random.rand(@map.length)


    @map.add(self)
  end
end

class DeadAnt < AntBase; end

class Ant < AntBase
  attr_accessor :radius, :density

  def think
    update_density

    if @dead_ant.nil?
      walk
      garbage_dead_ants if @random.rand < chance_to_garbage
    else
      walk
      drop_dead_ant if @random.rand < chance_to_drop
    end

    super
  end

  def garbaging
    !@dead_ant.nil?
  end

  private

  def chance_to_garbage
    1.0 - (@density.to_f / @higher_density.to_f)
  end

  def chance_to_drop
    @density.to_f / @higher_density.to_f
  end

  def update_density
    @higher_density ||= 0
    @density = @map.density(self)
    @higher_density = @density if @density > @higher_density
  end

  def update_anchor
    @anchor = @map.move_anchor(self)

    sum = 0

    for i in 0...3
      for j in 0...3
        sum += @anchor[i][j]
      end
    end

    for i in 0...3
      for j in 0...3
        @anchor[i][j] /= sum
      end
    end
  end

  def garbage_dead_ants
    @dead_ant = @map.get_dead_ant(ant: self)
  end

  def drop_dead_ant
    @dead_ant = nil if @map.drop_dead_ant(self, @dead_ant)
  end

  def walk
    @map.move(self) do
      @position_height = (@position_height + @random.rand(3) - 1) % @map.height
      @position_length = (@position_length + @random.rand(3) - 1) % @map.length
    end
  end

  def init
    @radius = 4
    @dead_ant = nil

    super
  end
end

class Map
  attr_accessor :length, :height

  def initialize(length: 50, height: 50, ants: 40, dead_ants: 500, ants_per_thread: 20)
    @length = length
    @height = height

    @specimens = ants
    @dead_specimens = dead_ants

    @ants_per_thread = ants_per_thread

    init
  end

  def move_anchor(ant)
    radius = ant.radius.to_i

    radius /= 2

    radius = 2 if radius < 1

    r = radius

    a = density_at(ant.position_length - r, ant.position_height - r, radius)
    b = density_at(ant.position_length - r, ant.position_height, radius)
    c = density_at(ant.position_length - r, ant.position_height + r, radius)
    d = density_at(ant.position_length, ant.position_height - r, radius)
    e = density_at(ant.position_length, ant.position_height + r, radius)
    f = density_at(ant.position_length + r, ant.position_height - r, radius)
    g = density_at(ant.position_length + r, ant.position_height, radius)
    h = density_at(ant.position_length + r, ant.position_height + r, radius)

    anchor = [
      [a, b, c],
      [d, 0, e],
      [f, g, h]
    ]


    anchor
  end

  def density(ant = nil)
    density_at(ant.position_length, ant.position_height, ant.radius)
  end

  def sync_map(length, height)
    @mutex_grid[length][height].synchronize do
      yield
    end
  end

  def drop_dead_ant(ant, dead_ant)
    return false if ant.nil? || dead_ant.nil? || @grid[ant.position_length][ant.position_height].count.zero?
    @grid[ant.position_length][ant.position_height] << dead_ant
  end

  def get_dead_ant(ant: nil)
    @grid[ant.position_length][ant.position_height].each do |object|
      next if object.class == Ant

      sync_map(ant.position_length, ant.position_height) do
        @grid[ant.position_length][ant.position_height].delete object
        return object
      end
    end

    nil
  end

  def move(ant = nil)
    remove(ant)
    yield
    add(ant)
  end

  def remove(ant = nil)
    unless ant.nil?
      @grid[ant.position_length][ant.position_height].delete ant
    end
  end

  def add(ant = nil)
    unless ant.nil?
      @grid[ant.position_length][ant.position_height] << ant
    end
  end

  def puts
    @grid.each do |line|
      line.each do |cell|
        case cell.last
        when Ant
          print '@'.blue if cell.last.garbaging
          print '@'.green unless cell.last.garbaging
        when DeadAnt
          print '#'.red
        else
          print ' '
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

  def density_at(x, y, radius)
    sum = 0.0
    dead_ants_count = 0
    for i in (-radius)..(radius)
      for j in (-radius)..(radius)
        sum += 1.0
        dead_ants_count += dead_ants_at((i + x) % @length, (j + y) % @height)
      end
    end

    dead_ants_count / sum
  end

  def dead_ants_at(i, j)
    @grid[i][j].select do |ant|
      ant.class == DeadAnt
    end.count
  end

  def init
    # init grid
    @grid ||= []

    for i in 0...@length
      @grid[i] ||= []

      for j in 0...@height
        @grid[i][j] ||= []
      end
    end

    # init mutex grid
    @mutex_grid ||= []
    for i in 0...@length
      @mutex_grid[i] ||= []

      for j in 0...@height
        @mutex_grid[i][j] ||= Mutex.new
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
map = Map.new(length: 50, height: 50, ants: 50, dead_ants: 500, ants_per_thread: 5)

frames = 300_000
step = 1

5_000_000.step(0, -step) do |n|
  system 'clear'
  puts "Frame #{n.to_f / frames}"
  puts Benchmark.measure { map.interate(step) }
  map.puts
end
